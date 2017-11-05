package demo

import (
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Client is a demo Client
type Client struct {
	ServerAddr string

	GRPCConn   *grpc.ClientConn
	GRPCClient DemoClient
}

func (c *Client) init() error {
	if c.GRPCClient != nil {
		return nil
	}
	var err error
	if c.GRPCConn == nil {
		c.GRPCConn, err = grpc.Dial(c.ServerAddr, grpc.WithInsecure())
		if err != nil {
			return err
		}
	}
	c.GRPCClient = NewDemoClient(c.GRPCConn)
	return nil
}

// Close close demo.Client
func (c *Client) Close() error {
	if c.GRPCConn != nil {
		return c.GRPCConn.Close()
	}
	return nil
}

// Echo call demo.Echo with gRPC
func (c *Client) Echo(ctx context.Context, req *EchoReq) (*EchoResp, error) {
	err := c.init()
	if err != nil {
		return nil, err
	}
	resp, err := c.GRPCClient.Echo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
