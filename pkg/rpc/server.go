package rpc

import (
	"fmt"
	"io"
	"net"
	"reflect"

	"github.com/Khighness/entry-task/pkg/log"
	"github.com/sirupsen/logrus"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-21

// Server struct
type Server struct {
	logger    *logrus.Logger
	addr      string
	functions map[string]reflect.Value
}

// NewServer create a new server
func NewServer(addr string) *Server {
	return &Server{
		logger:    log.NewLogger(logrus.WarnLevel, "", true),
		addr:      addr,                           // the net address of server
		functions: make(map[string]reflect.Value), // key: the name of func , value: reflect Value of function
	}
}

// Run start server
func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Infof("Listen at %s err: %v ", s.addr, err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Errorf("Accept client err: %v", err)
			continue
		} else {
			s.logger.Infof("Accept client: %s", conn.RemoteAddr())
		}

		go func() {
			srvTransport := NewTransport(conn)
			for {
				// read request from client
				req, err := srvTransport.Receive()
				if err != nil {
					if err != io.EOF {
						s.logger.Infof("Read request err: %v", err)
					}
					return
				}
				// get function by name
				f, ok := s.functions[req.Name]
				// if function requested does not exist
				if !ok {
					e := fmt.Sprintf("Func %s does not exist", req.Name)
					s.logger.Errorf(e)
					if err = srvTransport.Send(Data{Name: req.Name, Err: e}); err != nil {
						s.logger.Printf("Transport write err: %v", err)
					}
					continue
				}
				s.logger.Debugf("Call func: %s", req.Name)

				// un package function arguments
				inArgs := make([]reflect.Value, len(req.Args))
				for i := range req.Args {
					inArgs[i] = reflect.ValueOf(req.Args[i])
				}
				// invoke requested function
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
					s.logger.Errorf("Transport write err: %v", err)
				}
			}
		}()
	}
}

// Register register a function via name
func (s *Server) Register(name string, f interface{}) {
	if _, ok := s.functions[name]; ok {
		return
	}
	s.functions[name] = reflect.ValueOf(f)
	s.logger.Debugf("Register function: %v", name)
}
