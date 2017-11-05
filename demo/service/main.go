package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/anseres/anser/demo"
	"google.golang.org/grpc"
)

func main() {
	log.Printf("start service %s", serviceName)
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	demo.RegisterDemoServer(grpcServer, &demoServer{})
	grpcServer.Serve(lis)
	defer log.Printf("stop service %s", serviceName)
}
