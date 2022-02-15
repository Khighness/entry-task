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
	http.Handle("/", middleware.TimeMiddleWare(middleware.TokenMiddleWare(controller.Hello)))
	http.HandleFunc("/hello", controller.Hello)
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/register", controller.Register)
	log.Println("Start http server at port", web.HttpPort)
	err := http.ListenAndServe(web.HttpPort, nil)
	if err != nil {
		log.Println("Start http server failed")
		panic(err)
	}
}
