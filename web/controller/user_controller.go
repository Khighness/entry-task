package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Khighness/entry-task/pb"
	"github.com/Khighness/entry-task/pkg/rpc"
	"github.com/Khighness/entry-task/web/common"
	"github.com/Khighness/entry-task/web/config"
	"github.com/Khighness/entry-task/web/logging"
	"github.com/Khighness/entry-task/web/service"
	"github.com/Khighness/entry-task/web/util"
	"github.com/Khighness/entry-task/web/view"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// UserController 用户控制器
type UserController struct{}

// Ping Return Pong
func (userController *UserController) Ping(w http.ResponseWriter, r *http.Request) {
	view.HandleBizSuccess(w, "Pong")
}

// Register 用户注册
func (userController *UserController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		view.HandleMethodError(w, "Allowed Method: [POST]")
		return
	}
	var registerReq common.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		view.HandleRequestError(w, "Body should be json for registering data")
		return
	}

	var rpcRegister func(request pb.RegisterRequest) (pb.RegisterResponse, error)
	register := func(client *rpc.Client) {
		client.Call(pb.FuncRegister, &rpcRegister)
	}
	if err = service.Pool.Exec(register); err != nil {
		view.HandleErrorServerBusy(w)
		return
	}
	rpcRsp, _ := rpcRegister(pb.RegisterRequest{Username: registerReq.Username, Password: registerReq.Password})
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rpcRsp.Code, rpcRsp.Msg)
		return
	}
	view.HandleBizSuccess(w, nil)
}

// Login 用户登录
func (userController *UserController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		view.HandleMethodError(w, "Allowed Method: [POST]")
		return
	}
	var loginReq common.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		view.HandleRequestError(w, "Body should be json for logining data")
		return
	}

	var rpcLogin func(request pb.LoginRequest) (pb.LoginResponse, error)
	login := func(client *rpc.Client) {
		client.Call(pb.FuncLogin, &rpcLogin)
	}
	if err = service.Pool.Exec(login); err != nil {
		view.HandleErrorServerBusy(w)
		return
	}
	rpcRsp, _ := rpcLogin(pb.LoginRequest{Username: loginReq.Username, Password: loginReq.Password})
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rpcRsp.Code, rpcRsp.Msg)
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

	var rpcGetProfile func(request pb.GetProfileRequest) (pb.GetProfileResponse, error)
	getProfile := func(client *rpc.Client) {
		client.Call(pb.FuncGetProfile, &rpcGetProfile)
	}
	if err := service.Pool.Exec(getProfile); err != nil {
		view.HandleErrorServerBusy(w)
		return
	}
	rpcRsp, _ := rpcGetProfile(pb.GetProfileRequest{Token: r.Header.Get(common.HeaderTokenKey)})
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rpcRsp.Code, rpcRsp.Msg)
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
	var updateProfileReq common.UpdateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&updateProfileReq)
	if err != nil {
		view.HandleRequestError(w, "Body should be json for registering data")
		return
	}

	var rpcUpdateProfile func(request pb.UpdateProfileRequest) (pb.UpdateProfileResponse, error)
	updateProfile := func(client *rpc.Client) {
		client.Call(pb.FuncUpdateProfile, &rpcUpdateProfile)
	}
	if err = service.Pool.Exec(updateProfile); err != nil {
		view.HandleErrorServerBusy(w)
		return
	}
	rpcRsp, _ := rpcUpdateProfile(pb.UpdateProfileRequest{
		Token:    r.Header.Get(common.HeaderTokenKey),
		Username: updateProfileReq.Username,
	})
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rpcRsp.Code, rpcRsp.Msg)
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

	profilePicture := r.URL.Path[len(view.ShowAvatarUrl):]
	profilePicturePath := common.AvatarStoragePath + profilePicture
	_, err := os.Stat(profilePicturePath)
	if os.IsNotExist(err) {
		view.HandleBizError(w, profilePicture+" does not exist")
		logging.Log.Warn(profilePicturePath + " does not exist")
		return
	}
	file, _ := os.OpenFile(profilePicturePath, os.O_RDONLY, 0444)
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
	createFile, _ := os.OpenFile(common.AvatarStoragePath+avatarName, os.O_WRONLY|os.O_CREATE, 0766)
	defer createFile.Close()
	_, err := createFile.Write(fileStream)
	if err != nil {
		view.HandleBizError(w, "Upload profile picture failed")
		return
	}

	var rpcUpdateProfile func(request pb.UpdateProfileRequest) (pb.UpdateProfileResponse, error)
	updateProfile := func(client *rpc.Client) {
		client.Call(pb.FuncUpdateProfile, &rpcUpdateProfile)
	}
	if err = service.Pool.Exec(updateProfile); err != nil {
		view.HandleErrorServerBusy(w)
		return
	}
	rpcRsp, _ := rpcUpdateProfile(pb.UpdateProfileRequest{
		Token:          r.Header.Get(common.HeaderTokenKey),
		ProfilePicture: profilePicture,
	})
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rpcRsp.Code, rpcRsp.Msg)
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

	var rpcLogout func(request pb.LogoutRequest) (pb.LogoutResponse, error)
	logout := func(client *rpc.Client) {
		client.Call(pb.FuncLogout, &rpcLogout)
	}
	if err := service.Pool.Exec(logout); err != nil {
		view.HandleErrorServerBusy(w)
		return
	}
	rpcRsp, _ := rpcLogout(pb.LogoutRequest{Token: r.Header.Get(common.HeaderTokenKey)})
	if rpcRsp.Code != common.RpcSuccessCode {
		view.HandleErrorRpcResponse(w, rpcRsp.Code, rpcRsp.Msg)
		return
	}
	view.HandleBizSuccess(w, nil)
}
