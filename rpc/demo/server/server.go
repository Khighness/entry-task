package main

import (
	"encoding/gob"
	"entry/rpc"
	"entry/rpc/demo/public"
	"errors"
	"log"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-21

func queryUser(uid int64) (public.ResponseQueryUser, error) {
	db := make(map[int64]public.User)
	db[0] = public.User{Id:0, Name: "KHighness"}
	db[1] = public.User{Id:1, Name: "FlowerK"}
	if u, ok := db[uid]; ok {
		return public.ResponseQueryUser{User: u, Msg: "success"}, nil
	}
	return public.ResponseQueryUser{User: public.User{}, Msg: "fail"}, errors.New("uid is not in database")
}

func main() {
	gob.Register(public.ResponseQueryUser{})

	addr := "127.0.0.1:30000"
	srv := rpc.NewServer(addr)
	srv.Register("queryUser", queryUser)
	log.Printf("Server is running at %v\n", addr)
	go srv.Run()

	for {
	}
}

