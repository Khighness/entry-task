package rpc

import (
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-21

// Server struct
type Server struct {
	addr      string
	functions map[string]reflect.Value
}

// NewServer create a new server
func NewServer(addr string) *Server {
	return &Server{
		addr:      addr,
		functions: make(map[string]reflect.Value),
	}
}

// Run start server
func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("Listen at %s err: %v \n", s.addr, err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept client err: %v\n", err)
			continue
		} else {
			log.Printf("Accept client %s\n", conn.RemoteAddr())
		}

		go func() {
			srvTransport := NewTransport(conn)
			for {
				// read request from client
				req, err := srvTransport.Receive()
				if err != nil {
					if err != io.EOF {
						log.Printf("Read request err: %v\n", err)
					}
					return
				}
				// get method by name
				f, ok := s.functions[req.Name]
				// if method requested does not exist
				if !ok {
					e := fmt.Sprintf("Func %s does not exist", req.Name)
					log.Printf(e)
					if err = srvTransport.Send(Data{Name: req.Name, Err: e}); err != nil {
						log.Printf("Transport write err: %v\n", err)
					}
					continue
				}
				log.Printf("Func %s is called\n", req.Name)

				// unPackage request arguments
				inArgs := make([]reflect.Value, len(req.Args))
				for i := range req.Args {
					inArgs[i] = reflect.ValueOf(req.Args[i])
				}
				// invoke requested method
				out := f.Call(inArgs)
				// package response arguments
				outArgs := make([]interface{}, len(out)-1)
				for i := range req.Args {
					outArgs[i] = out[i].Interface()
				}
				// package error argument
				var e string
				if _, ok := out[len(out)-1].Interface().(error); !ok {
					e = ""
				} else {
					e = out[len(out)-1].Interface().(error).Error()
				}
				// send response to client
				err = srvTransport.Send(Data{Name: req.Name, Args: outArgs, Err: e})
				if err != nil {
					log.Printf("Transport write err: %v\n", err)
				}
			}
		}()
	}
}

// Register a method via name
func (s *Server) Register(name string, f interface{}) {
	if _, ok := s.functions[name]; ok {
		return
	}
	s.functions[name] = reflect.ValueOf(f)
}
