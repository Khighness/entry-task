package main

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/Khighness/entry-task/pkg/rpc"
	"github.com/Khighness/entry-task/pkg/rpc/demo/public"
)

// @Author KHighness
// @Email  zikang.chen@shopee.com
// @Since 2022-02-20

var QueryUser func(int64) (public.ResponseQueryUser, error)

func main() {
	gob.Register(public.ResponseQueryUser{})

	conn, err := net.Dial("tcp", "0.0.0.0:30000")
	if err != nil {
		log.Fatalln(err)
	}
	client := rpc.NewClient(conn)
	client.Call("queryUser", &QueryUser)
	response, err := QueryUser(1)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", response)
}
