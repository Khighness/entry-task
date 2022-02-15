## entry-task

### 中间件部署

部署MySQL


```shell
$ mkdir -p /Users/zikang.chen/Docker/mysql/data /Users/zikang.chen/Docker/mysql/conf
$ docker run --name mysql -d -p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=KAG1823 mysql:8.0.20
$ docker cp mysql:/etc/mysql/my.cnf /Users/zikang.chen/Docker/mysql/conf
$ vim /Users/zikang.chen/Docker/mysql/conf/my.cnf
# ADD
[mysqld]
character-set-server=utf8
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
$ ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'KAG1823';

```



部署Redis


```shell
$ mkdir -p /Users/zikang.chen/Docker/redis/data /Users/zikang.chen/Docker/redis/conf
$ cd /Users/zikang.chen/Docker/redis/conf
$ touch redis.conf
$ cat << EOF >>/Users/zikang.chen/Docker/redis/conf/redis.conf
port 6379
#bind 0.0.0.0
daemonize no
protected-mode no
requirepass KAG1823
loglevel verbose

maxmemory-policy volatile-ttl
slowlog-log-slower-than 2000
maxclients 512
timeout 1800

dir ./
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
auto-aof-rewrite-min-size 128mb
dbfilename dump.rdb
save 900 1

cluster-enabled no
cluster-config-file nodes.conf
cluster-node-timeout 5000
cluster-announce-port 6379
cluster-announce-bus-port 16379
EOF

$ docker run -d -p 6379:6379 --name redis \
-v /Users/zikang.chen/Docker/redis/data:/data \
-v /Users/zikang.chen/Docker/redis/conf/redis.conf:/etc/redis/redis.conf \
redis:6.2.6 \
--requirepass "KAG1823" 
```
