package router

import (
	"entry/web/common"
	"entry/web/controller"
	"entry/web/logging"
	"entry/web/middleware"
	"entry/web/view"
	"net/http"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// Start 启动web server
func Start() {
	userController := controller.UserController{}
	http.HandleFunc(view.RegisterUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(userController.Register)))
	http.HandleFunc(view.LoginUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(userController.Login)))
	http.HandleFunc(view.GetProfileUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(middleware.TokenMiddleWare(userController.GetProfile))))
	http.HandleFunc(view.UpdateProfileUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(middleware.TokenMiddleWare(userController.UpdateProfile))))
	http.HandleFunc(view.ShowAvatarUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(userController.ShowAvatar)))
	http.HandleFunc(view.UploadAvatarUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(middleware.TokenMiddleWare(userController.UploadAvatar))))
	http.HandleFunc(view.LogoutUrl, middleware.CorsMiddleWare(middleware.TimeMiddleWare(userController.Logout)))
	logging.Log.Infof("Web server is serving at [%s]", common.HttpServerAddr)
	err := http.ListenAndServe(common.HttpServerAddr, nil)
	if err != nil {
		logging.Log.Fatalf("Failed to start web server at [%s]", common.HttpServerAddr)
	}
}
