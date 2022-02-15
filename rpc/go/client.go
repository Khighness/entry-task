package main

import (
	"log"
	"net/rpc"
)

// @Author KHighness
// @Update 2022-02-16

func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("Client connect to server failed:", err)
	}

	var reply string
	err = client.Call("HelloService.Hello", "KHighness", &reply)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(reply)
}