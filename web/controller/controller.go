package controller

import (
	"entry/tcp/util"
	web "entry/web/common"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
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

// Hello 测试
func Hello(w http.ResponseWriter, r *http.Request) {
	var answer = fmt.Sprintf("{\"status\":\"ok\"， \"data\":\"hello world\"}")
	w.Write([]byte(answer))
}

// Login 用户登录
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == web.Get {
		t, _ := template.ParseFiles("./public/template/login.gtpl")
		t.Execute(w, nil)
	} else if r.Method == web.Post {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])

		// TODO 生成sessionId，存在cookie中
		cookie := http.Cookie{
			Name:     web.CookieSessionKey,
			Value:    util.GenerateSessionId(),
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24),
			HttpOnly: false,
		}
		http.SetCookie(w, &cookie)
		fmt.Fprintf(w, "ok")
	}
}

// Register 用户注册
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == web.Get {
		t, _ := template.ParseFiles("./public/template/register.gtpl")
		t.Execute(w, nil)
	} else if r.Method == web.Post {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		fmt.Fprintf(w, "ok")
	}
}

// UploadAvatar 上传头像
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method == web.Get {
		t, _ := template.ParseFiles("./public/template/profile.gtpl")
		t.Execute(w, nil)
	} else if r.Method == web.Post {
		r.ParseForm()
		uploadFile, header, err := r.FormFile("profile_picture")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		if !(strings.HasSuffix(header.Filename, ".jpg") || strings.HasSuffix(header.Filename, ".png") || strings.HasSuffix(header.Filename, ".jpeg")) {
			w.Write([]byte("请上传jpg/png/jpeg格式文件"))
			return
		}
		avatarName := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
		// TODO avatarName 存入数据库
		defer uploadFile.Close()
		openFile, err := os.OpenFile(web.AvatarPath + avatarName, os.O_WRONLY | os.O_CREATE, 0777)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer openFile.Close()
		_, err = io.Copy(openFile, uploadFile)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		log.Println("用户上传头像: ", avatarName)
		w.Write([]byte("ok"))
	}
}

// ShowAvatar 显示头像
func ShowAvatar(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/avatar/")
	file, err := os.Open(web.AvatarPath + path)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(buf)
}
