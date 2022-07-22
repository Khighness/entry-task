package main

import (
	"encoding/gob"
	"errors"
	"log"

	"github.com/Khighness/entry-task/rpc"
	"github.com/Khighness/entry-task/rpc/demo/public"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-21

type userService struct{}

func (u userService) queryUser(uid int64) (public.ResponseQueryUser, error) {
	db := make(map[int64]public.User)
	db[0] = public.User{Id: 0, Name: "KHighness"}
	db[1] = public.User{Id: 1, Name: "FlowerK"}
	if u, ok := db[uid]; ok {
		return public.ResponseQueryUser{User: u, Msg: "success"}, nil
	}
	return public.ResponseQueryUser{User: public.User{}, Msg: "fail"}, errors.New("uid is not in database")
}

func main() {
	gob.Register(public.ResponseQueryUser{})

	addr := "127.0.0.1:30000"
	srv := rpc.NewServer(addr)
	service := userService{}
	srv.Register("queryUser", service.queryUser)
	log.Printf("Server is running at %v\n", addr)
	srv.Run()
}
