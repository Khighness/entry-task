package server

import (
	"fmt"
	"github.com/Khighness/entry-task/tcp/model"

	"github.com/Khighness/entry-task/pb"
	"github.com/Khighness/entry-task/pkg/rpc"
	"github.com/Khighness/entry-task/tcp/cache"
	"github.com/Khighness/entry-task/tcp/config"
	"github.com/Khighness/entry-task/tcp/logging"
	"github.com/Khighness/entry-task/tcp/mapper"
	"github.com/Khighness/entry-task/tcp/service"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

// Start 启动TCP server
func Start() {
	// 创建RPC Server
	serverCfg := config.AppCfg.Server
	serverAddr := fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port)
	server := rpc.NewServer(serverAddr)

	// 注册Function
	userCache := cache.NewUserCache(cache.RedisClient)
	userMapper := mapper.NewUserMapper(model.DB)
	userService := service.NewUserService(userMapper, userCache)
	server.Register(pb.FuncRegister, userService.Register)
	server.Register(pb.FuncLogin, userService.Login)
	server.Register(pb.FuncCheckToken, userService.CheckToken)
	server.Register(pb.FuncGetProfile, userService.GetProfile)
	server.Register(pb.FuncUpdateProfile, userService.UpdateProfile)
	server.Register(pb.FuncLogout, userService.Logout)

	// 启动RPC Server
	logging.Log.Printf("RPC tcp server is serving at [%s]", serverAddr)
	server.Run()
}
