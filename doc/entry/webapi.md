## entry-task web api document



🚀 状态码：10000为SUCCESS，其他均为ERROR。

| code  | message                                                  |
| ----- | -------------------------------------------------------- |
| 10000 | SUCCESS                                                  |
| 20000 | ERROR                                                    |
| 20001 | Server is busy, please try again later                   |
| 20002 | RPC failed or timeout                                    |
| 30001 | 用户名长度不得小于3                                      |
| 30002 | 用户名长度不得大于18                                     |
| 30003 | 用户名已存在，请换个试试                                 |
| 30004 | 密码长度不得小于6                                        |
| 30005 | 密码长度不得大于20                                       |
| 30006 | 密码长度较弱，最少需要包含数字/字母/特殊符号中的以上两种 |
| 30007 | 用户名错误                                               |
| 30008 | 密码错误                                                 |
| 30009 | 令牌非法                                                 |
| 30010 | 登录状态已过期                                           |
| 40001 | 操作数据库失败                                           |



### 1. 用户注册

请求方式：POST

请求路径：/register

请求数据：`application/json`

| 字段     | 类型   | 描述   |
| -------- | ------ | ------ |
| username | string | 用户名 |
| password | string | 密码   |

返回数据：`application/json`

| 字段    | 类型   | 描述                       |
| ------- | ------ | -------------------------- |
| code    | int    | 状态码                     |
| message | string | 信息，注册失败时为提示信息 |
| data    | string | null                       |

请求示例：

```shell
$ curl -X POST \
     -d '{"username":"Khighness","password":"czk911"}' \
     'http://127.0.0.1:10000/register'
{
	"code": 10000,
	"message": "SUCCESS",
	"data": null
}
```



### 2. 用户登录

请求方式：POST

请求路径：/login

请求数据：`application/json`

| 字段     | 类型   | 描述   |
| -------- | ------ | ------ |
| username | string | 用户名 |
| password | string | 密码   |

返回数据：`application/json`

| 字段    | 类型   | 描述           |
| ------- | ------ | -------------- |
| code    | int    | 状态码         |
| message | string | 信息           |
| data    | json   | 令牌和用户信息 |

请求示例：

```shell
$ curl -X POST \
     -d '{"username":"Khighness","password":"czk911"}' \
     'http://127.0.0.1:10000/login'
{
	"code": 10000,
	"message": "SUCCESS",
	"data": {
		"token": "1b17da6fb96ee7bf9bf1ca11a1d68703",
		"user": {
			"id": 1,
			"username": "Khighness",
			"profilePicture": "http://127.0.0.1:10000/avatar/show/khighness.jpg"
		}
	}
}
```



### 3. 个人信息

请求方式：GET

请求路径：/user/profile

请求头部：`Authorization`

返回数据：`application/json`

| 字段    | 类型   | 描述     |
| ------- | ------ | -------- |
| code    | int    | 状态码   |
| message | string | 信息     |
| data    | json   | 用户信息 |

请求示例：
```shell
$ curl -X GET \
     -H 'Authorization:720df5d0c0e9649597f54b531a1e348d'	\
     'http://127.0.0.1:10000/user/profile'
{
	"code": 10000,
	"message": "SUCCESS",
	"data": {
		"id": 1,
		"username": "KHighness",
		"profilePicture": "http://127.0.0.1:10000/avatar/show/khighness.jpg"
	}
}
```



### 4. 更新信息

请求方式：PUT

请求路径：/user/update

请求头部：`Authorization`

请求数据：`application/json`

| 字段     | 类型   | 描述                             |
| -------- | ------ | -------------------------------- |
| username | string | 用户名                           |

返回数据：`application/json`

| 字段    | 类型   | 描述   |
| ------- | ------ | ------ |
| code    | int    | 状态码 |
| message | string | 信息   |
| data    | string | null   |

请求示例：

```shell
$ curl -X PUT \
     -H 'Authorization:1b17da6fb96ee7bf9bf1ca11a1d68703' \
     -d '{"username":"Khighness1"}' \
     'http://127.0.0.1:10000/user/update'
{
	"code": 10000,
	"message": "SUCCESS",
	"data": null
}
```



### 5. 展示图片

请求方式：GET

请求路径：/avatar/show/${picture}

返回数据：`image/jpg` / `image/png` / `image/jpeg`

请求示例：

```shell
$ curl -X GET \
     -o '/Users/zikang.chen/Pictures/output.jpg' \
     'http://127.0.0.1:10000/avatar/show/khighness.jpg'
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 30100    0 30100    0     0  14.3M      0 --:--:-- --:--:-- --:--:-- 14.3M
```



### 6. 更新头像

请求方式：POST

请求路径：/avatar/upload

请求头部：`Authorization`

请求数据：`multi/form-data`

| 字段            | 类型 | 描述     |
| --------------- | ---- | -------- |
| profile_picture | file | 头像文件 |

返回数据：

请求示例：

```shell
$ curl -X POST \
     -H 'Authorization:1b17da6fb96ee7bf9bf1ca11a1d68703' \
     -F 'profile_picture=@/Users/zikang.chen/Pictures/Khighness.jpg' \
     'http://127.0.0.1:10000/avatar/upload'
{
	"code": 10000,
	"message": "SUCCESS",
	"data": null
}
```



### 7. 退出登录

请求方式：get

请求路径：/logout

请求头部：`Authorization`

返回数据：`application/json`

| 字段    | 类型   | 描述   |
| ------- | ------ | ------ |
| code    | int    | 状态码 |
| message | string | 信息   |
| data    | string | null   |

请求示例：

```shell
$ curl -X GET \
     -H 'Authorization:1b17da6fb96ee7bf9bf1ca11a1d68703' \
     'http://127.0.0.1:10000/logout'
{
	"code": 10000,
	"message": "SUCCESS",
	"data": null
}
```
