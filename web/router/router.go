package router

import (
	"entry/web/api"
	"entry/web/common"
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
	http.HandleFunc("/", api.Index)
	http.HandleFunc(view.RegisterUrl, middleware.TimeMiddleWare(api.Register))
	http.HandleFunc(view.LoginUrl, middleware.TimeMiddleWare(api.Login))
	http.HandleFunc(view.ProfileUrl, middleware.TimeMiddleWare(middleware.TokenMiddleWare(api.GetProfile)))
	http.HandleFunc(view.AvatarUrl, middleware.TimeMiddleWare(api.ShowAvatar))
	http.HandleFunc(view.UpdateUrl, middleware.TimeMiddleWare(middleware.TokenMiddleWare(api.UpdateInfo)))
	log.Printf("Web server is serving at [%s]\n", common.HttpAddr)
	err := http.ListenAndServe(common.HttpAddr, nil)
	if err != nil {
		log.Fatalf("Failed to start web server at [%s]\n", common.HttpAddr)
	}
}
