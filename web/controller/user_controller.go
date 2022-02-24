package controller

import (
	"context"
	"encoding/json"
	"entry/pb"
	"entry/web/common"
	"entry/web/grpc"
	"entry/web/view"
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

type UserController struct{}

// Register 用户注册
func (userController *UserController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		view.HandleMethodError(w, "Allowed Method: [GET]")
		return
	}
	var registerRequest common.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		view.HandleRequestError(w, "Body should be json for registering data")
		return
	}

	permission, err := grpc.Pool.Achieve(context.Background())
	defer grpc.Pool.Release(permission, context.Background())
	if err != nil {
		view.HandlerBizError(w, "Server is busy, please try again later")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rpcRsp, err := permission.RpcCli.Register(ctx, &pb.RegisterRequest{
		Username: registerRequest.Username,
		Password: registerRequest.Password,
	})

	if err != nil {
		view.HandlerBizError(w, "RPC failed or timeout")
		return
	}
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandlerRpcErrResponse(w, rpcRsp.Code, rpcRsp.Msg)
		return
	}
	view.HandleBizSuccess(w, nil)
}

// Login 用户登录
// TODO 防止XSRF攻击
func (userController *UserController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		view.HandleMethodError(w, "Allowed Method: [GET]")
		return
	}
	var loginRequest common.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		view.HandleRequestError(w, "Body should be json for logining data")
		return
	}

	permission, err := grpc.Pool.Achieve(context.Background())
	defer grpc.Pool.Release(permission, context.Background())
	if err != nil {
		view.HandlerBizError(w, "Server is busy, please try again later")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rpcRsp, err := permission.RpcCli.Login(ctx, &pb.LoginRequest{
		Username: loginRequest.Username,
		Password: loginRequest.Password,
	})

	if err != nil {
		view.HandlerBizError(w, "RPC failed or timeout")
		return
	}
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandlerRpcErrResponse(w, rpcRsp.Code, rpcRsp.Msg)
		return
	}
	view.HandleBizSuccess(w, common.LoginResponse{
		Token: rpcRsp.Token,
		User: common.UserInfo{
			Id:             rpcRsp.User.Id,
			Username:       rpcRsp.User.Username,
			ProfilePicture: rpcRsp.User.ProfilePicture,
		},
	})
}

// GetProfile 获取信息
func (userController *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		view.HandleMethodError(w, "Allowed Method: [GET]")
		return
	}

	permission, err := grpc.Pool.Achieve(context.Background())
	if err != nil {
		view.HandlerBizError(w, "Server is busy, please try again later")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rpcRsp, err := permission.RpcCli.GetProfile(ctx, &pb.GetProfileRequest{
		Token: r.Header.Get(common.HeaderTokenKey),
	})
	go grpc.Pool.Release(permission, context.Background())

	if err != nil {
		view.HandlerBizError(w, "RPC failed or timeout")
		return
	}
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandlerRpcErrResponse(w, rpcRsp.Code, rpcRsp.Msg)
		return
	}
	view.HandleBizSuccess(w, common.UserInfo{
		Id:             rpcRsp.User.Id,
		Username:       rpcRsp.User.Username,
		ProfilePicture: rpcRsp.User.ProfilePicture,
	})
}

// UpdateProfile 更新信息
func (userController *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		view.HandleMethodError(w, "Allowed Method: [PUT]")
		return
	}
	var updateProfileRequest common.UpdateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&updateProfileRequest)
	if err != nil {
		view.HandleRequestError(w, "Body should be json for registering data")
		return
	}

	permission, err := grpc.Pool.Achieve(context.Background())
	defer grpc.Pool.Release(permission, context.Background())
	if err != nil {
		view.HandlerBizError(w, "Server is busy, please try again later")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rpcRsp, err := permission.RpcCli.UpdateProfile(ctx, &pb.UpdateProfileRequest{
		Token:    r.Header.Get(common.HeaderTokenKey),
		Username: updateProfileRequest.Username,
	})

	if err != nil {
		view.HandlerBizError(w, "RPC failed or timeout")
		return
	}
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandlerRpcErrResponse(w, rpcRsp.Code, rpcRsp.Msg)
		return
	}
	view.HandleBizSuccess(w, nil)
}

// ShowAvatar 显示头像
func (userController *UserController) ShowAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		view.HandleMethodError(w, "Allowed Method: [GET]")
		return
	}
	profilePicture := strings.TrimLeft(r.URL.Path, view.ShowAvatarUrl)
	_, err := os.Stat(common.AvatarStoragePath + profilePicture)
	if os.IsNotExist(err) {
		view.HandleBizSuccess(w, profilePicture+" does not exist")
		return
	}
	file, _ := os.OpenFile(common.AvatarStoragePath+profilePicture, os.O_RDONLY, 0444)
	defer file.Close()
	buf, _ := ioutil.ReadAll(file)
	_, _ = w.Write(buf)
}

// UploadAvatar 上传头像
func (userController *UserController) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		view.HandleMethodError(w, "Allowed Method: [POST]")
		return
	}

	uploadFile, header, _ := r.FormFile("profile_picture")
	if uploadFile == nil {
		view.HandleRequestError(w, "Picture file cannot be null")
		return
	}
	defer uploadFile.Close()
	if !(strings.HasSuffix(header.Filename, ".jpg") || strings.HasSuffix(header.Filename, ".png") || strings.HasSuffix(header.Filename, ".jpeg")) {
		view.HandleRequestError(w, "Please upload jpg/png/jpeg file as profile picture")
		return
	}

	avatarName := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
	profilePicture := fmt.Sprintf("http://%s%s%s", common.HttpServerAddr, view.ShowAvatarUrl, avatarName)
	createFile, err := os.OpenFile(common.AvatarStoragePath+avatarName, os.O_WRONLY|os.O_CREATE, 0666)
	defer createFile.Close()
	_, err = io.Copy(createFile, uploadFile)
	if err != nil {
		view.HandlerBizError(w, "Upload profile picture failed")
		return
	}

	permission, err := grpc.Pool.Achieve(context.Background())
	defer grpc.Pool.Release(permission, context.Background())
	if err != nil {
		view.HandlerBizError(w, "Server is busy, please try again later")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rpcRsp, err := permission.RpcCli.UpdateProfile(ctx, &pb.UpdateProfileRequest{
		Token:          r.Header.Get(common.HeaderTokenKey),
		ProfilePicture: profilePicture,
	})

	if err != nil {
		view.HandlerBizError(w, "RPC failed or timeout")
		return
	}
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandlerRpcErrResponse(w, rpcRsp.Code, rpcRsp.Msg)
		return
	}
	view.HandleBizSuccess(w, nil)
}

// Logout 退出登录
func (userController *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		view.HandleMethodError(w, "Allowed Method: [GET]")
		return
	}

	permission, err := grpc.Pool.Achieve(context.Background())
	defer grpc.Pool.Release(permission, context.Background())
	if err != nil {
		view.HandlerBizError(w, "Server is busy, please try again later")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rpcRsp, err := permission.RpcCli.Logout(ctx, &pb.LogoutRequest{
		Token: r.Header.Get(common.HeaderTokenKey),
	})

	if err != nil {
		view.HandlerBizError(w, "RPC failed or timeout")
		return
	}
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandlerRpcErrResponse(w, rpcRsp.Code, rpcRsp.Msg)
		return
	}
	view.HandleBizSuccess(w, nil)
}
