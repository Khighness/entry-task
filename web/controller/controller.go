package controller

import (
	web "entry/web/common"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

func Index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "ok")
}

// 测试
func Hello(w http.ResponseWriter, r *http.Request) {
	var answer = fmt.Sprintf("{\"status\":\"ok\"， \"data\":\"hello world\"}")
	w.Write([]byte(answer))
}

// 登录
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == web.Get {
		t, _ := template.ParseFiles("./public/login.gtpl")
		t.Execute(w, nil)
	} else if r.Method == web.Post {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		fmt.Fprintf(w, "ok")
	}
}

// 注册
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == web.Get {
		t, _ := template.ParseFiles("./public/register.gtpl")
		t.Execute(w, nil)
	} else if r.Method == web.Post {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		fmt.Fprintf(w, "ok")
	}
}

// 上传头像
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method == web.Get {
		t, _ := template.ParseFiles("./public/profile.gtpl")
		t.Execute(w, nil)
	} else if r.Method == web.Post {
		r.ParseForm()
		//file, header, err := r.FormFile("profile_picture")

	}
}

// 显示头像
func ShowAvatar(w http.ResponseWriter, r *http.Request) {

}

// 错误处理
func handleError(err error, w http.ResponseWriter) {
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
