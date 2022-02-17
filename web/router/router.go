package router

import (
	"entry/web/api"
	"entry/web/common"
	"entry/web/middleware"
	"log"
	"net/http"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// Start 启动http server
func Start() {
	http.HandleFunc("/hello", api.Hello)
	http.HandleFunc("/login", middleware.TimeMiddleWare(api.Login))
	http.HandleFunc("/register", middleware.TimeMiddleWare(api.Register))
	http.HandleFunc("/update", middleware.TimeMiddleWare(api.UpdateInfo))
	http.HandleFunc("/upload", middleware.TimeMiddleWare(api.UploadAvatar))
	http.HandleFunc("/avatar/", middleware.TimeMiddleWare(api.ShowAvatar))
	log.Printf("Http server is serving at [%s]\n", common.HttpAddr)
	err := http.ListenAndServe(common.HttpAddr, nil)
	if err != nil {
		log.Fatalf("Failed to start http server at [%s]\n", common.HttpAddr)
	}
}
