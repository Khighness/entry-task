## entry-task

![golang](https://img.shields.io/badge/go-1.17-blue?style=for-the-badge&logo=go)



### 项目结构

```
entry-task
    ├─doc              文档
    ├─pb               proto
    ├─rpc              todo
    ├─tcp              TCP服务器
    │  ├─cache         redis缓存
    │  ├─common        初始化
    │  ├─conf          tcp配置
    │  ├─mapper        持久化
    │  ├─model         服务器
    │  ├─service       rpc接口
    │  └─util          工具
    └─web              WEB服务器
        ├─api          web接口
        ├─common       初始化
        ├─conf         web配置
        ├─grpc         rpc调用
        ├─middleware   web中间件
        ├─public       静态文件
        ├─router       web路由
        └─view         模板解析
```





### 环境部署

部署MySQL


```shell
$ mkdir -p /Users/zikang.chen/Docker/mysql/data /Users/zikang.chen/Docker/mysql/doc
$ docker run --name mysql -d -p 3306:3306 \
-e MYSQL_ROOT_PASSWORD=KAG1823 mysql:8.0.20
$ docker cp mysql:/etc/mysql/my.cnf /Users/zikang.chen/Docker/mysql/doc
$ vim /Users/zikang.chen/Docker/mysql/doc/my.cnf
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
-v /Users/zikang.chen/Docker/mysql/doc/my.cnf:/etc/mysql/my.cnf \
-v /Users/zikang.chen/Docker/mysql/data:/var/lib/mysql \
--restart=on-failure:3 \
mysql:8.0.20
$ docker exec -it mysql bash
$ mysql -u root -p KAG1823
$ ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'KAG1823';

```



部署Redis


```shell
$ mkdir -p /Users/zikang.chen/Docker/redis/data /Users/zikang.chen/Docker/redis/doc
$ cd /Users/zikang.chen/Docker/redis/doc
$ touch redis.doc
$ cat << EOF >>/Users/zikang.chen/Docker/redis/doc/redis.doc
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
-v /Users/zikang.chen/Docker/redis/doc/redis.doc:/etc/redis/redis.doc \
redis:6.2.6 \
--requirepass "KAG1823" 
```



### 快速运行

1. 下载依赖

```shell
$ go mod tidy
```

2. 启动tcp server

```shell
$ go run tcp/main.go
```

3. 启动web server

```shell
$ go run web/main.go
```



> 效果预览

<table>
  <tr>
    <td><a href="http://127.0.0.1:10000/login">登录</td>
    <td><a href="http://127.0.0.1:10000/profile">个人</td>
  </tr>
  <tr>
     <td width="50%" align="top"><img src="./doc/images/login.png"/></td>
     <td width="50%" align="top"><img src="./doc/images/profile.png"/></td>
  </tr>
</table>
