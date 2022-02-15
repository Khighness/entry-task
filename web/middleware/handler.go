package middleware

import (
	"log"
	"net/http"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// MiddleHandler 封装接口
type MiddleHandler func(next http.HandlerFunc) http.HandlerFunc

// TimeMiddleWare 耗时计算与日志打印
func TimeMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 处理业务耗时
		startTime := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(startTime)

		log.Printf("[%v] url:%v, method:%v, time:%v \n", r.RemoteAddr, r.URL.Path, r.Method, timeElapsed)
	})
}

// TokenMiddleWare 用户Token认证
func TokenMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("rpc tcp server，进行token验证")
		// 从cookie中取出sessionId

		// tcp server校验sessionId是否合法

		// 继续处理业务
		next.ServeHTTP(w, r)
	})
}
