package server

import (
	"fmt"

	"github.com/Khighness/entry-task/public"
	"github.com/Khighness/entry-task/rpc"
	"github.com/Khighness/entry-task/tcp/cache"
	"github.com/Khighness/entry-task/tcp/config"
	"github.com/Khighness/entry-task/tcp/logging"
	"github.com/Khighness/entry-task/tcp/mapper"
	"github.com/Khighness/entry-task/tcp/service"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

// Start 启动tcp server
func Start() {
	// 创建RPC Server
	serverCfg := config.AppCfg.Server
	serverAddr := fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port)
	server := rpc.NewServer(serverAddr)

	// 注册Function
	userService := service.NewUserService(&mapper.UserMapper{}, &cache.UserCache{})
	server.Register(public.FuncRegister, userService.Register)
	server.Register(public.FuncLogin, userService.Login)
	server.Register(public.FuncCheckToken, userService.CheckToken)
	server.Register(public.FuncGetProfile, userService.GetProfile)
	server.Register(public.FuncUpdateProfile, userService.UpdateProfile)
	server.Register(public.FuncLogout, userService.Logout)

	// 启动RPC Server
	logging.Log.Printf("GRPC tcp server is serving at [%s]", serverAddr)
	server.Run()
}
