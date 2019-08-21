package rpc

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	"github.com/sah4ez/grpc2nats/pkg/types"
	"google.golang.org/grpc"
)

type Client struct {
	nc *nats.Conn
}

func (c *Client) Generate(ctx context.Context, in *types.GenerateRequest, opts ...grpc.CallOption) (*types.GenerateResponse, error) {
	resp := &types.GenerateResponse{}
	b, err := proto.Marshal(in)
	if err != nil {
		return resp, err
	}
	msg, err := c.nc.Request(subjGenerate, b, time.Second*5)
	if err != nil {
		return resp, err
	}
	err = proto.Unmarshal(msg.Data, resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *Client) run() {
}

func NewClient(nc *nats.Conn) *Client {
	return &Client{
		nc: nc,
	}
}
