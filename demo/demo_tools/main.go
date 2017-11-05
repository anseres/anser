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
	client, err := demo.NewClient("127.0.0.1:21000")
	if err != nil {
		panic(err)
	}
	req := &demo.EchoReq{
		Message: "hello anser",
	}
	resp, err := client.Echo(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Echo: req=%s\nresp=%s\n", req.String(), resp.String())
}
