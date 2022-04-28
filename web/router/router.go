package router

import (
	"fmt"
	"github.com/Khighness/entry-task/web/config"
	"github.com/Khighness/entry-task/web/controller"
	"github.com/Khighness/entry-task/web/logging"
	"github.com/Khighness/entry-task/web/middleware"
	"github.com/Khighness/entry-task/web/view"
	"net/http"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// Start 启动web server
func Start() {
	userController := controller.UserController{}
	http.HandleFunc(view.PingUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(userController.Ping)))
	http.HandleFunc(view.RegisterUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(userController.Register)))
	http.HandleFunc(view.LoginUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(userController.Login)))
	http.HandleFunc(view.GetProfileUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(middleware.TokenMiddleWare(userController.GetProfile))))
	http.HandleFunc(view.UpdateProfileUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(middleware.TokenMiddleWare(userController.UpdateProfile))))
	http.HandleFunc(view.ShowAvatarUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(userController.ShowAvatar)))
	http.HandleFunc(view.UploadAvatarUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(middleware.TokenMiddleWare(userController.UploadAvatar))))
	http.HandleFunc(view.LogoutUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(userController.Logout)))

	serverCfg := config.AppCfg.Server
	serverAddr := fmt.Sprintf("%s:%d", serverCfg.Host, serverCfg.Port)
	logging.Log.Infof("Web server is serving at [%s]", serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		logging.Log.Fatalf("Failed to start web server, error: %s", err)
	}
}
