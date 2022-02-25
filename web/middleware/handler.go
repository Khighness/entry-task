package middleware

import (
	"context"
	"entry/pb"
	"entry/web/common"
	"entry/web/grpc"
	"entry/web/logging"
	"entry/web/view"
	"net/http"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// MiddleHandler 封装接口
type MiddleHandler func(next http.HandlerFunc) http.HandlerFunc

// CorsMiddleWare 处理跨域
func CorsMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 允许访问所有域
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许header值
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		// 允许携带cookie
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		// 允许请求方法
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// 允许返回任意格式数据
		w.Header().Set("content-type", "*")

		// 跨域第一次OPTIONS请求，直接放行
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	}

}

// TimeMiddleWare 日志打印
func TimeMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 处理业务耗时
		startTime := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(startTime)

		logging.Log.Infof("[IP:%v] url:%v, method:%v, time:%v", r.RemoteAddr, r.URL.Path, r.Method, timeElapsed)
	})
}

// TokenMiddleWare Token认证
func TokenMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从header中取出token
		token := r.Header.Get(common.HeaderTokenKey)
		if token == "" {
			view.HandleBizError(w, "Authorization failed")
			return
		}
		logging.Log.Debugf("[verify token] token: %s", token)

		// 校验token是否合法
		permission, err := grpc.Pool.Achieve(context.Background())
		defer grpc.Pool.Release(permission, context.Background())
		if err != nil {
			view.HandleErrorServerBusy(w)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		rpcRsp, err := permission.RpcCli.CheckToken(ctx, &pb.CheckTokenRequest{Token: token})
		if err != nil {
			view.HandleErrorRpcRequest(w, permission)
			return
		}

		// 认证失败
		if rpcRsp.Code != common.RpcSuccessCode {
			view.HandleErrorRpcResponse(w, rpcRsp.Code, rpcRsp.Msg, permission)
			return
		}

		// 认证成功，继续处理业务
		next.ServeHTTP(w, r)
	})
}
