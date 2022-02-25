package main

import (
	"encoding/gob"
	"entry/rpc/demo/public"
)

// @Author KHighness
// @Update 2022-02-20

var QueryUser func(int64) (public.ResponseQueryUser, error)

func main() {
	gob.Register(public.ResponseQueryUser{})

}
