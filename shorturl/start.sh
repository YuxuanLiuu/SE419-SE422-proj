## mysql
kubectl create -f shorturl-mysql-pv.yaml
sleep 1
kubectl create -f shorturl-mysql-deployment.yaml
## redis
kubectl create -f redis.yaml
## shorturl server
kubectl create -f shorturl-go.yaml -f expose-service.yaml
