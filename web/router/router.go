package router

import (
	web "entry/web/common"
	"entry/web/controller"
	"entry/web/middleware"
	"log"
	"net/http"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// NewRouter 启动http服务器
func NewRouter() {
	http.HandleFunc("/hello", controller.Hello)
	http.HandleFunc("/login", middleware.TimeMiddleWare(controller.Login))
	http.HandleFunc("/register", middleware.TimeMiddleWare(controller.Register))
	http.HandleFunc("/upload", middleware.TimeMiddleWare(controller.UploadAvatar))
	http.HandleFunc("/avatar/", middleware.TimeMiddleWare(controller.ShowAvatar))
	log.Println("Start http server at port", web.HttpPort)
	err := http.ListenAndServe(web.HttpPort, nil)
	if err != nil {
		log.Println("Start http server failed")
		panic(err)
	}
}
