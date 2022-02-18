package api

import (
	"context"
	"entry/pb"
	"entry/web/common"
	"entry/web/grpc"
	"entry/web/view"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// Index 测试
func Index(w http.ResponseWriter, r *http.Request) {
	view.HandleError(w, "Error", "机房失火断电<br>节点故障宕机<br>服务熔断降级", "Sign In", view.LoginUrl)
}

// Register 用户注册
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == common.Get {
		view.DirectRegister(w)
	} else if r.Method == common.Post {
		_ = r.ParseForm()
		username := strings.Join(r.Form["username"], "")
		password := strings.Join(r.Form["password"], "")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, err := grpc.Client.Register(ctx, &pb.RegisterRequest{
			Username: username,
			Password: password,
		})

		if err != nil {
			view.HandleError(w, common.DefaultErrorType, common.DefaultErrorMessage, "Sign Up", view.RegisterUrl)
			return
		}
		if response.Code != common.RpcSuccessCode {
			view.HandleError(w, "注册失败", response.Msg, "Sign Up", view.RegisterUrl)
			return
		}

		view.HandleSuccess(w, "注册成功", fmt.Sprintf("亲爱的用户%s，感谢您的支持", username), "Sign In", view.LoginUrl)
	}
}

// Login 用户登录
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == common.Get {
		view.DirectLogin(w)
	} else if r.Method == common.Post {
		_ = r.ParseForm()
		username := strings.Join(r.Form["username"], "")
		password := strings.Join(r.Form["password"], "")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, err := grpc.Client.Login(ctx, &pb.LoginRequest{
			Username: username,
			Password: password,
		})

		if err != nil {
			view.HandleError(w, common.DefaultErrorType, common.DefaultErrorMessage, "Sign In", view.LoginUrl)
			return
		}
		if response.Code != common.RpcSuccessCode {
			view.HandleError(w, "登录失败", response.Msg, "Sign In", view.LoginUrl)
			return
		}

		http.SetCookie(w, generateCookie(common.CookieTokenKey, response.SessionId))
		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

// GetProfile 获取信息
func GetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCookie(r)
	if err != nil {
		view.HandleError(w, common.CookieErrorType, common.CookieErrorMessage, "Sign In", view.LoginUrl)
		return
	}
	view.DirectProfile(w, common.UserInfo{
		Id:             user.Id,
		Username:       user.Username,
		ProfilePicture: user.ProfilePicture,
	})
}

// ShowAvatar 显示头像
func ShowAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method == common.Get {
		file, err := os.Open(common.FileStoragePath + r.URL.Path)
		if err != nil {
			view.HandleError(w, "显示头像", "没找到哎", "Sign In", view.LoginUrl)
			return
		}
		defer file.Close()
		buf, _ := ioutil.ReadAll(file)
		_, _ = w.Write(buf)
	}
}

// UpdateInfo 更新信息
func UpdateInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == common.Get {
		user, err := getUserFromCookie(r)
		if err != nil {
			view.HandleError(w, "更新失败", err.Error(), "Update Profile", view.UpdateUrl)
			return
		}
		view.DirectUpdate(w, *user)
	} else if r.Method == common.Post {
		_ = r.ParseMultipartForm(1024)
		cookie, _ := r.Cookie(common.CookieTokenKey)
		username := strings.Join(r.Form["username"], "")
		profilePicture := ""

		// 检查用户是否上传了头像
		uploadFile, header, _ := r.FormFile("profile_picture")
		if uploadFile != nil {

			if !(strings.HasSuffix(header.Filename, ".jpg") || strings.HasSuffix(header.Filename, ".png") || strings.HasSuffix(header.Filename, ".jpeg")) {
				view.HandleError(w, "更新失败", "请上传jpg/png/jpeg格式文件作为头像", "Update Profile", view.UpdateUrl)
				return
			}

			avatarName := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
			profilePicture = fmt.Sprintf("http://%s/%s/%s", common.HttpAddr, common.RelativeAvatarPath, avatarName)

			// 存储文件
			defer uploadFile.Close()
			createFile, err := os.OpenFile(common.FileStoragePath+common.RelativeAvatarPath+avatarName, os.O_WRONLY|os.O_CREATE, 0777)
			defer createFile.Close()
			_, err = io.Copy(createFile, uploadFile)
			if err != nil {
				view.HandleError(w, "更新失败", "上传头像失败", "Update Profile", view.UpdateUrl)
				return
			}
		}

		// 更新数据库
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, err := grpc.Client.UpdateProfile(ctx, &pb.UpdateProfileRequest{
			SessionId:      cookie.Value,
			Username:       username,
			ProfilePicture: profilePicture,
		})

		if err != nil {
			view.HandleError(w, common.DefaultErrorType, common.DefaultErrorMessage, "Update", view.UpdateUrl)
			return
		}
		if response.Code != common.RpcSuccessCode {
			view.HandleError(w, "更新失败", response.Msg, "Update Profile", view.UpdateUrl)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

// getUserFromCookie 根据cookie获取用户信息
func getUserFromCookie(r *http.Request) (user *common.UserInfo, err error) {
	cookie, err := r.Cookie(common.CookieTokenKey)
	if err != nil {
		return nil, errors.New(common.CookieErrorMessage)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := grpc.Client.GetProfile(ctx, &pb.GetProfileRequest{SessionId: cookie.Value})
	if err != nil || response.Code != common.RpcSuccessCode {
		return nil, errors.New(common.DefaultErrorMessage)
	}

	return &common.UserInfo{
		Id:             response.User.Id,
		Username:       response.User.Username,
		ProfilePicture: response.User.ProfilePicture,
	}, nil
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
