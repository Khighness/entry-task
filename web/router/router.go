package router

import (
	"entry/web/common"
	"entry/web/controller"
	"entry/web/middleware"
	"entry/web/view"
	"log"
	"net/http"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// Start 启动web server
func Start() {
	userController := controller.UserController{}
	http.HandleFunc("/", userController.Index)
	http.HandleFunc(view.RegisterUrl, middleware.TimeMiddleWare(userController.Register))
	http.HandleFunc(view.LoginUrl, middleware.TimeMiddleWare(userController.Login))
	http.HandleFunc(view.ProfileUrl, middleware.TimeMiddleWare(middleware.TokenMiddleWare(userController.GetProfile)))
	http.HandleFunc(view.AvatarUrl, middleware.TimeMiddleWare(userController.ShowAvatar))
	http.HandleFunc(view.UpdateUrl, middleware.TimeMiddleWare(middleware.TokenMiddleWare(userController.UpdateInfo)))
	http.HandleFunc(view.LogoutUrl, middleware.TimeMiddleWare(userController.Logout))
	log.Printf("Web server is serving at [%s]\n", common.HttpServerAddr)
	err := http.ListenAndServe(common.HttpServerAddr, nil)
	if err != nil {
		log.Fatalf("Failed to start web server at [%s]\n", common.HttpServerAddr)
	}
}
