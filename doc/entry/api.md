## 接口文档



### 1. 用户注册

（1）前端页面

请求方式：GET

请求路径：/register

请求结果：html

请求示例：

```shell
curl -X GET \
    'http://127.0.0.1:10000/register'
```

（2）后端服务

请求方式：POST

请求路径：/register

请求结果：html

请求参数：form

| 字段     | 类型   | 描述   |
| -------- | ------ | ------ |
| username | string | 用户名 |
| password | string | 密码   |

请求示例：

```shell
curl -X POST \
     -d 'username=testname&password=test1234' \
     'http://127.0.0.1:10000/register'
```



### 2. 用户登录

（1）前端页面

请求方式：GET

请求路径：/login

请求结果：html

请求示例：

```shell
curl -X GET \
     'http://127.0.0.1:10000/login'
```

（2）后端服务

请求方式：POST

请求路径：/login

请求结果：redirect

请求参数：`application/x-www-form-urlencoded`

| 字段     | 类型   | 描述   |
| -------- | ------ | ------ |
| username | string | 用户名 |
| password | string | 密码   |

请求示例：

```shell
curl -X POST \
     -d 'username=Khighness&password=czk911' \
     'http://127.0.0.1:10000/login'
```



### 3. 个人信息

请求方式：GET

请求路径：/profile

请求结果：html

请求示例：
```shell
curl -X GET \
     --cookie 'sessionId=16656C030BA870FDC835523FD3317040'	\
     'http://127.0.0.1:10000/profile'
```



### 4. 更新信息

（1）前端页面

请求方式：GET

请求路径：/update

请求结果：html

请求示例：

```shell
curl -X GET \
     'http://127.0.0.1:10000/update'
```

（2）后端服务

请求方式：GET

请求路径：/profile

请求结果：redirect

请求参数：`multipart/form-data`

| 字段            | 类型   | 描述                             |
| --------------- | ------ | -------------------------------- |
| profile_picture | File   | 用户头像，仅限[jpg/png/jpeg]格式 |
| password        | string | 用户名                           |

请求示例：

```shell
curl -X POST \
     -F 'profile_picture=@/Users/zikang.chen/Pictures/Khighness.jpg' \
     --cookie 'sessionId=16656C030BA870FDC835523FD3317040' \
     'http://127.0.0.1:10000/update?username=RubbishK'
```



### 5. 显示图片

请求方式：GET

请求路径：/avatar/${picture}

请求结果：binary

```shell
curl -X GET \
     -o '/Users/zikang.chen/Pictures/output.jpg' \
     'http://127.0.0.1:10000/avatar/default.jpg'
```

