package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
)

const (
	userName = "root"
	password = "13549621996yx"
	// password = "shortUrl"
	dbName = "shorturl"
)

//Db数据库连接池
var DB *sql.DB
var redisDB *redis.Client

//注意方法名大写，就是public
func InitDB() {
	ipMysql := os.Getenv("SHORT_URL_MYSQL_IP")
	portMysql := os.Getenv("SHORT_URL_MYSQL_PORT")
	ipRedis := os.Getenv("SHORT_URL_REDIS_IP")
	portRedis := os.Getenv("SHORT_URL_REDIS_PORT")
	if strings.Count(ipMysql, "")-1 == 0 {
		// ipMysql = "shorturl-mysql"
		ipMysql = "localhost"
	}
	if strings.Count(portMysql, "")-1 == 0 {
		portMysql = "3306"
	}

	if strings.Count(ipRedis, "")-1 == 0 {
		//ipRedis = "shorturl-redis"
		ipRedis = "localhost"
	}
	if strings.Count(portRedis, "")-1 == 0 {
		portRedis = "6379"
	}
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ipMysql, ":", portMysql, ")/", dbName, "?charset=utf8"}, "")
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetMaxOpenConns(1000)
	DB.SetConnMaxLifetime(500 * time.Minute)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("open mysql database fail")
		return
	}
	fmt.Println("connnect mysql success")

	redisDB = redis.NewClient(&redis.Options{
		Addr:     ipRedis + ":" + portRedis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisDB.Ping().Result()
	if err != nil {
		fmt.Println("open redis database fail")
		return
	}
	fmt.Println("connnect redis success")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func convertChar2Int(char byte) int {
	ord := int(char)
	if char >= 'A' && char <= 'Z' {
		return ord - int('A')
	} else if char >= 'a' && char <= 'z' {
		return ord - 'a' + 26
	} else if char >= '0' && char <= '9' {
		return ord - '0' + 52
	} else if char == '+' {
		return 62
	} else if char == '-' {
		return 63
	} else {
		return 0
	}
}

func convertInt2Char(i int8) byte {
	if i >= 0 && i <= 25 {
		return 'A' + byte(i)
	} else if i >= 26 && i <= 51 {
		return 'a' + byte(i-26)
	} else if i >= 52 && i <= 61 {
		return '0' + byte(i-52)
	} else if i == 62 {
		return '+'
	} else if i == 63 {
		return '-'
	} else {
		return 0
	}
}

func convertString2Int(str string) int {
	sum := 0
	for i := 0; i < len(str); i++ {
		sum += convertChar2Int(str[i]) << (i * 6)
	}
	return sum
}

func convertInt2String(sum int64) string {
	var bytes [6]byte
	for i := 0; i < 6; i++ {
		bytes[i] = convertInt2Char(int8(sum % 64))
		sum /= 64
	}
	return string(bytes[:])
}

func shortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fullURL, err := redisDB.Get(r.URL.Path[1:7]).Result()
		if err == nil {
			println("here redis")
			fmt.Fprintf(w, fullURL)
		} else {
			println("here mysql")
			path := strconv.Itoa(convertString2Int(r.URL.Path[1:7]))
			stmt, err := DB.Prepare("SELECT `long` from short where id=?")
			check(err)
			var long string
			err = stmt.QueryRow(path).Scan(&long)
			if err != nil {
				fmt.Fprintf(w, "failed to acquire long url")
			} else {
				redisDB.Set(r.URL.Path[1:7], long, time.Hour)
				fmt.Fprintf(w, long)
			}
			stmt.Close()
		}
	} else if r.Method == "POST" {
		stmt, err := DB.Prepare("Insert into short values(null, ?)")
		check(err)
		url := r.FormValue("url")

		res, err := stmt.Exec(url)
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Fprintf(w, "fail")
		} else {
			fmt.Fprintf(w, r.Host+r.URL.String()+convertInt2String(id))
		}
		stmt.Close()
	}
	r.Body.Close()
}

func main() {
	InitDB()
	http.HandleFunc("/", shortURL)
	errInfo := http.ListenAndServe(":8081", nil)
	if errInfo != nil {
		log.Fatal("ListenAndServe:", errInfo)
	}
	defer DB.Close()
	defer redisDB.Close()
}
