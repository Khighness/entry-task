package hello

import "net/rpc"

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

const HelloServiceName = "path/to/pkg.HelloService"

type HelloServiceInterface interface {
	Hello(request string, reply *string) error
}

func RegisterHelloService(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}
