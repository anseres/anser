package main

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/anseres/anser/demo"
	context "golang.org/x/net/context"
)

func main() {
	grpc.WithInsecure()
	ctx := context.Background()
	client := demo.Client{}
	client.ServerAddr = "127.0.0.1:21000"
	req := &demo.EchoReq{
		Message: "hello anser",
	}
	resp, err := client.Echo(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Echo: req=%s\nresp=%s\n", req.String(), resp.String())
}
