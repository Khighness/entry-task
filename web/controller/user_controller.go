package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Khighness/entry-task/pb"
	"github.com/Khighness/entry-task/web/common"
	"github.com/Khighness/entry-task/web/config"
	"github.com/Khighness/entry-task/web/grpc"
	"github.com/Khighness/entry-task/web/logging"
	"github.com/Khighness/entry-task/web/util"
	"github.com/Khighness/entry-task/web/view"
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

// Ping Return Pong
func (userController *UserController) Ping(w http.ResponseWriter, r *http.Request) {
	view.HandleBizSuccess(w, "Pong")
}

// Register 用户注册
func (userController *UserController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		view.HandleMethodError(w, "Allowed Method: [GET]")
		return
	}
	var registerReq common.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		view.HandleRequestError(w, "Body should be json for registering data")
		return
	}

	register := func(cli pb.UserServiceClient) (interface{}, error) {
		return cli.Register(context.Background(), &pb.RegisterRequest{Username: registerReq.Username, Password: registerReq.Password})
	}
	rpcRsp, err := grpc.GP.Exec(register)
	if err != nil {
		view.HandleErrorRpcRequest(w)
		return
	}
	rsp, _ := rpcRsp.(*pb.RegisterResponse)
	if rsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rsp.Code, rsp.Msg)
		return
	}
	view.HandleBizSuccess(w, nil)
}

// Login 用户登录
func (userController *UserController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		view.HandleMethodError(w, "Allowed Method: [GET]")
		return
	}
	var loginReq common.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		view.HandleRequestError(w, "Body should be json for logining data")
		return
	}

	login := func(cli pb.UserServiceClient) (interface{}, error) {
		return cli.Login(context.Background(), &pb.LoginRequest{Username: loginReq.Username, Password: loginReq.Password})
	}
	rpcRsp, err := grpc.GP.Exec(login)
	if err != nil {
		view.HandleErrorRpcRequest(w)
		return
	}
	rsp, _ := rpcRsp.(*pb.LoginResponse)
	if rsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rsp.Code, rsp.Msg)
		return
	}
	view.HandleBizSuccess(w, common.LoginResponse{
		Token: rsp.Token,
		User: common.UserInfo{
			Id:             rsp.User.Id,
			Username:       rsp.User.Username,
			ProfilePicture: rsp.User.ProfilePicture,
		},
	})
}

// GetProfile 获取信息
func (userController *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		view.HandleMethodError(w, "Allowed Method: [GET]")
		return
	}

	getProfile := func(cli pb.UserServiceClient) (interface{}, error) {
		return cli.GetProfile(context.Background(), &pb.GetProfileRequest{Token: r.Header.Get(common.HeaderTokenKey)})
	}
	rpcRsp, err := grpc.GP.Exec(getProfile)
	if err != nil {
		view.HandleErrorRpcRequest(w)
		return
	}
	rsp, _ := rpcRsp.(*pb.GetProfileResponse)
	if rsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rsp.Code, rsp.Msg)
		return
	}
	view.HandleBizSuccess(w, common.UserInfo{
		Id:             rsp.User.Id,
		Username:       rsp.User.Username,
		ProfilePicture: rsp.User.ProfilePicture,
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

	updateProfile := func(cli pb.UserServiceClient) (interface{}, error) {
		return cli.UpdateProfile(context.Background(), &pb.UpdateProfileRequest{
			Token:    r.Header.Get(common.HeaderTokenKey),
			Username: updateProfileRequest.Username,
		})
	}
	rpcRsp, err := grpc.GP.Exec(updateProfile)
	if err != nil {
		view.HandleErrorRpcRequest(w)
		return
	}
	rsp, _ := rpcRsp.(*pb.UpdateProfileResponse)
	if rsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rsp.Code, rsp.Msg)
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
	profilePictirePath := common.AvatarStoragePath + profilePicture
	_, err := os.Stat(profilePictirePath)
	if os.IsNotExist(err) {
		view.HandleBizError(w, profilePicture+" does not exist")
		logging.Log.Warn(profilePictirePath + " does not exist")
		return
	}
	file, _ := os.OpenFile(profilePictirePath, os.O_RDONLY, 0444)
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
	if !(strings.HasSuffix(header.Filename, "jpg") || strings.HasSuffix(header.Filename, "png") || strings.HasSuffix(header.Filename, "jpeg")) {
		view.HandleRequestError(w, "Please upload jpg/png/jpeg file as profile picture")
		return
	}
	fileStream, _ := io.ReadAll(uploadFile)
	fileType := util.GetFileType(fileStream[:10])
	if !(fileType == "jpg" || fileType == "png" || fileType == "jpeg") {
		view.HandleRequestError(w, "The suffix of the uploaded file does not match the content of the file")
		return
	}
	avatarName := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)

	serverCfg := config.AppCfg.Server
	serverAddr := fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port)
	profilePicture := fmt.Sprintf("http://%s%s%s", serverAddr, view.ShowAvatarUrl, avatarName)
	createFile, err := os.OpenFile(common.AvatarStoragePath+avatarName, os.O_WRONLY|os.O_CREATE, 0766)
	defer createFile.Close()
	_, err = createFile.Write(fileStream)
	if err != nil {
		view.HandleBizError(w, "Upload profile picture failed")
		return
	}

	updateProfile := func(cli pb.UserServiceClient) (interface{}, error) {
		return cli.UpdateProfile(context.Background(), &pb.UpdateProfileRequest{
			Token:          r.Header.Get(common.HeaderTokenKey),
			ProfilePicture: profilePicture,
		})
	}
	rpcRsp, err := grpc.GP.Exec(updateProfile)
	if err != nil {
		view.HandleErrorRpcRequest(w)
		return
	}
	rsp, _ := rpcRsp.(*pb.UpdateProfileResponse)
	if rsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rsp.Code, rsp.Msg)
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

	logout := func(cli pb.UserServiceClient) (interface{}, error) {
		return cli.Logout(context.Background(), &pb.LogoutRequest{Token: r.Header.Get(common.HeaderTokenKey)})
	}
	rpcRsp, err := grpc.GP.Exec(logout)
	if err != nil {
		view.HandleErrorRpcRequest(w)
		return
	}
	rsp, _ := rpcRsp.(*pb.LogoutResponse)
	if rsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rsp.Code, rsp.Msg)
		return
	}
	view.HandleBizSuccess(w, nil)
}
