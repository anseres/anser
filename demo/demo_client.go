package demo

import (
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Client is a Demo Client
type Client struct {
	ServerAddr string

	GRPCConn   *grpc.ClientConn
	GRPCClient DemoClient
}

func NewClient(addr string) (*Client, error) {
    client := &Client{}
    client.ServerAddr = addr
    err := client.init()
    if err != nil {
        return nil, err
    }
    return client, nil
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

// Close close Demo.Client
func (c *Client) Close() error {
	if c.GRPCConn != nil {
		return c.GRPCConn.Close()
	}
	return nil
}

// Echo call Demo.Echo with gRPC
func (c *Client) Echo(ctx context.Context, req *EchoReq) (*EchoResp, error) {
	resp, err := c.GRPCClient.Echo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}