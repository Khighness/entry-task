package public

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-21

// User struct
type User struct {
	Id   int64
	Name string
}

// ResponseQueryUser struct
type ResponseQueryUser struct {
	User
	Msg string
}