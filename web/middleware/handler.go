package middleware

import (
	"context"
	"entry/pb"
	"entry/web/common"
	"entry/web/grpc"
	"entry/web/view"
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

		log.Printf("[%v] url:%v, method:%v, time:%v\n", r.RemoteAddr, r.URL.Path, r.Method, timeElapsed)
	})
}

// TokenMiddleWare 用户Token认证
func TokenMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从cookie中取出sessionId
		cookie, err := r.Cookie(common.CookieTokenKey)
		if err != nil {
			view.HandleError(w, common.CookieErrorType, common.CookieErrorMessage, "Sign In", view.LoginUrl)
			return
		}

		// 校验sessionId是否合法
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		permission, err := grpc.Pool.Achieve(ctx)
		if err != nil {
			log.Println(err)
			view.HandleError(w, common.DefaultErrorType, common.DefaultErrorMessage, "Sign Up", view.RegisterUrl)
			return
		}
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		response, err := permission.RpcCli.CheckToken(ctx, &pb.CheckTokenRequest{Token: cookie.Value})
		if err != nil {
			view.HandleError(w, common.CookieErrorType, common.CookieErrorMessage, "Sign In", view.LoginUrl)
			return
		}
		go grpc.Pool.Release(permission.RpcCli, context.Background())

		if response.Code == common.RpcSuccessCode {
			// 认证成功，继续处理业务
			next.ServeHTTP(w, r)
		} else {
			// 认证失败，重定向到登录
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	})
}
