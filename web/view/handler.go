package view

import (
	"entry/web/common"
	"net/http"
	"text/template"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

// DirectRegister 注册页面
func DirectRegister(w http.ResponseWriter) {
	t, _ := template.ParseGlob(templatePath + RegisterUrl + templateSuffix)
	_ = t.Execute(w, nil)
}

// DirectLogin 登陆页面
func DirectLogin(w http.ResponseWriter) {
	t, _ := template.ParseGlob(templatePath + LoginUrl + templateSuffix)
	_ = t.Execute(w, nil)
}

// DirectProfile 个人页面
func DirectProfile(w http.ResponseWriter, user common.UserInfo) {
	t, _ := template.ParseGlob(templatePath + ProfileUrl + templateSuffix)
	_ = t.Execute(w, user)
}

// DirectUpdate 更新页面
func DirectUpdate(w http.ResponseWriter, user common.UserInfo) {
	t, _ := template.ParseGlob(templatePath + UpdateUrl + templateSuffix)
	_ = t.Execute(w, user)
}

// HandleSuccess 成功页面
func HandleSuccess(w http.ResponseWriter, sucType, message, returnTip, returnUrl string) {
	t, _ := template.ParseGlob("./web/public/template/success.html")
	_ = t.Execute(w, common.SuccessMsg{
		SucType:   sucType,
		Message:   message,
		ReturnTip: returnTip,
		ReturnUrl: returnUrl,
	})
}

// HandleError 错误页面
func HandleError(w http.ResponseWriter, errType, message, returnTip, returnUrl string) {
	t, _ := template.ParseGlob("./web/public/template/error.html")
	_ = t.Execute(w, common.ErrorMsg{
		ErrType:   errType,
		Message:   message,
		ReturnTip: returnTip,
		ReturnUrl: returnUrl,
	})
}
