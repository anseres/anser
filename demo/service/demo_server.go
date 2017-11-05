package main

import (
	"log"

	"github.com/anseres/anser/demo"
	context "golang.org/x/net/context"
)

var serviceName = "echo"
var port = 21000

type demoServer struct {
}

func (s *demoServer) Echo(ctx context.Context, req *demo.EchoReq) (*demo.EchoResp, error) {
	log.Printf("echo: %s", req.Message)
	return &demo.EchoResp{Result: req.Message}, nil
}
