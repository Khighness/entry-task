package api

import (
	"context"
	"entry/pb"
	"entry/web/common"
	"entry/web/grpc"
	"entry/web/view"
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

// Hello 测试
func Hello(w http.ResponseWriter, r *http.Request) {
	view.HandleError(w, "Error", "机房失火断电<br>节点故障宕机<br>服务熔断降级")
	//ShowProfile(w, r, common.UserInfo{
	//	Id:             1,
	//	Username:       "Khighness",
	//	ProfilePicture: "http://127.0.0.1:10000/avatar/Khighness.jpg",
	//})
}

// Register 用户注册
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == common.Get {
		view.DirectRegister(w)
	} else if r.Method == common.Post {
		_ = r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

// Login 用户登录
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == common.Get {
		view.DirectLogin(w)
	} else if r.Method == common.Post {
		_ = r.ParseForm()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, err := grpc.Client.Login(ctx, &pb.LoginRequest{
			Username: strings.Join(r.Form["username"], ""),
			Password: strings.Join(r.Form["password"], ""),
		})

		if err != nil {
			view.HandleError(w, "", "")
			return
		}
		if response.Code != common.RpcSuccessCode {
			view.HandleError(w, "登录失败", response.Msg)
		}

		http.SetCookie(w, generateCookie(common.CookieTokenKey, response.SessionId))
	}
}

// UpdateInfo 更新信息
func UpdateInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == common.Get {
		view.DirectUpdate(w, common.UserInfo{
			Id:             0,
			Username:       "KHighness",
			ProfilePicture: "",
		})
	} else if r.Method == common.Post {
		_ = r.ParseMultipartForm(1024)

		username := strings.Join(r.Form["username"], "")
		avatarName := ""
		log.Println("username:", username)

		// 检查用户是否上传了头像
		uploadFile, header, err := r.FormFile("profile_picture")
		if err != nil {
			return
		}
		if uploadFile != nil {
			if !(strings.HasSuffix(header.Filename, ".jpg") || strings.HasSuffix(header.Filename, ".png") || strings.HasSuffix(header.Filename, ".jpeg")) {
				view.HandleError(w, "更新失败", "请上传jpg/png/jpeg格式文件作为头像")
				return
			}
			avatarName = fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
		}

		// 存储头像
		defer uploadFile.Close()
		createFile, err := os.OpenFile(common.FileStoragePath+"avatar/"+avatarName, os.O_WRONLY|os.O_CREATE, 0777)
		defer createFile.Close()
		_, err = io.Copy(createFile, uploadFile)
		if err != nil {
			view.HandleError(w, "更新失败", "上传头像失败")
			return
		}

		// 更新数据库

	}
}

// UploadAvatar 上传头像
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method == common.Get {
		t, _ := template.ParseGlob("./public/template/profile.html")
		t.Execute(w, nil)
	} else if r.Method == common.Post {
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
		openFile, err := os.OpenFile(common.FileStoragePath+avatarName, os.O_WRONLY|os.O_CREATE, 0777)
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
	if r.Method == common.Get {
		file, err := os.Open(common.FileStoragePath + r.URL.Path)
		if err != nil {
			view.HandleError(w, "显示头像", "Not Found")
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				view.HandleError(w, "", "")
			}
		}(file)
		buf, _ := ioutil.ReadAll(file)
		_, _ = w.Write(buf)
	} else {
		view.HandleError(w, "显示头像", "Error Method")
	}
}

// getUserFromCookie 根据cookie获取用户信息
func getUserFromCookie(r *http.Request) {
	//cookie, err := r.Cookie(common.CookieTokenKey)
}

// generateCookie 生成cookie
func generateCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(common.CookieTokenTimeout),
		HttpOnly: false,
	}
}
