## entry task deploy document



### 🐬 MySQL


```shell
$ mkdir -p /Users/zikang.chen/Docker/mysql/data /Users/zikang.chen/Docker/mysql/conf
$ docker run --name mysql -d -p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=KAG1823 mysql:8.0.20
$ docker cp mysql:/etc/mysql/my.cnf /Users/zikang.chen/Docker/mysql/conf
$ vim /Users/zikang.chen/Docker/mysql/conf/my.cnf
# ADD
[mysqld]
character-set-server=utf8
max_connections=30000
[client]
default-character-set=utf8
[mysql]
default-character-set=utf8

$ docker stop mysql && docker rm mysql
$ docker run --name mysql \
-d -p 3306:3306  \
-e MYSQL_ROOT_PASSWORD=KAG1823 \
-v /Users/zikang.chen/Docker/mysql/conf/my.cnf:/etc/mysql/my.cnf \
-v /Users/zikang.chen/Docker/mysql/data:/var/lib/mysql \
--restart=on-failure:3 \
mysql:8.0.20
$ docker exec -it mysql bash
$ mysql -u root -p KAG1823
$ ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'KANG1823';

```



### 💠 Redis


```shell
$ mkdir -p /Users/zikang.chen/Docker/redis/data /Users/zikang.chen/Docker/redis/conf
$ cd /Users/zikang.chen/Docker/redis/conf
$ touch redis.conf
$ cat << EOF >>/Users/zikang.chen/Docker/redis/conf/redis.conf
port 6379
#bind 0.0.0.0
daemonize no
protected-mode no
requirepass KANG1823
loglevel notice

maxmemory-policy volatile-ttl
slowlog-log-slower-than 2000
maxclients 30000
timeout 3600

dir /usr/local/redis/data/
appendonly yes
appendfilename "appendonly.aof"
appendfsync no
auto-aof-rewrite-min-size 128mb
dbfilename dump.rdb
save 900 1
EOF

$ docker run -d -p 6379:6379 --name redis \
-v /Users/zikang.chen/Docker/redis/data:/data \
-v /Users/zikang.chen/Docker/redis/conf/redis.conf:/etc/redis/redis.conf \
redis:6.2.6 \
--requirepass "KANG1823" 
```



### 🚀 Start

1. 导入脚本

```
./doc/mysql/db.sql
```

2. 下载依赖

```shell
$ go mod tidy
```

3. 启动tcp server

```shell
$ go run tcp/main.go
```

4. 启动web server

```shell
$ go run web/main.go
```

5. 启动vue

```shell
$ cd front
$ npm install
$ npm run serve
```

TODO 二进制文件
