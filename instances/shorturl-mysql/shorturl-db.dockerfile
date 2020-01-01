FROM mysql:8

ENV MYSQL_ROOT_PASSWORD shortUrl
# ENV MYSQL_ALLOW_EMPTY_PASSWORD=true
# ENV MYSQL_DATABASE=socksdb

COPY ./dump.sql /docker-entrypoint-initdb.d/

