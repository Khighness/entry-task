package grpc

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Khighness/entry-task/pb"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-03-08

func Test(t *testing.T) {
	connector := GrpcConnector{GrpcServerAddr: "127.0.0.1:20000"}
	grpcPool := NewGrpcPool(connector, &GrpcPoolConfig{
		MaxOpenCount: 10,
		MaxIdleCount: 5,
		MaxLifeTime:  10 * time.Second,
		MaxIdleTime:  5 * time.Second,
	})

	go func() {
		req := &pb.LoginRequest{
			Username: "Khighness",
			Password: "czk911",
		}
		f := func(cli pb.UserServiceClient) (interface{}, error) {
			return cli.Login(context.Background(), req)
		}
		rsp, err := grpcPool.Exec(f)
		if err != nil {
			fmt.Println(err)
		} else {
			valueOf := reflect.ValueOf(rsp)
			response := valueOf.Interface().(*pb.LoginResponse)
			fmt.Printf("%+v\n", response)
		}
	}()

	time.Sleep(3 * time.Second)
	fmt.Printf("%+v\n", grpcPool.Stat())
}
