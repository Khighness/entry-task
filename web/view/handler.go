package view

import (
	"entry/web/common"
	"net/http"
	"text/template"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

// DirectLogin 登陆页面
func DirectLogin(w http.ResponseWriter) {
	t, _ := template.ParseFiles("./web/public/template/login.html")
	_ = t.Execute(w, nil)
}

// DirectRegister 注册页面
func DirectRegister(w http.ResponseWriter) {
	t, _ := template.ParseGlob("./web/public/template/register.html")
	_ = t.Execute(w, nil)
}

// DirectProfile 个人页面
func DirectProfile(w http.ResponseWriter, user common.UserInfo) {
	t, _ := template.ParseGlob("./web/public/template/profile.html")
	_ = t.Execute(w, user)
}

// DirectUpdate 更新页面
func DirectUpdate(w http.ResponseWriter, user common.UserInfo) {
	t, _ := template.ParseGlob("./web/public/template/update.html")
	_ = t.Execute(w, user)
}

// HandleError 错误页面
func HandleError(w http.ResponseWriter, errType, message string) {
	if errType == "" {
		errType = common.DefaultErrorType
	}
	if message == "" {
		message = common.DefaultErrorMessage
	}
	t, _ := template.ParseGlob("./web/public/template/error.html")
	_ = t.Execute(w, common.ErrorMsg{
		ErrType: errType,
		Message: message,
	})
}
