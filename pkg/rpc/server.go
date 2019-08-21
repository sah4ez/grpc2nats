package rpc

import (
	"context"
	"crypto/md5"
	"fmt"
	"hash"
	"sync"

	"github.com/golang/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	"github.com/sah4ez/grpc2nats/pkg/types"
	"github.com/satori/go.uuid"
)

type Server struct {
	nc   *nats.Conn
	one  *sync.Once
	subs map[string]*nats.Subscription
	h    hash.Hash
}

func (s *Server) Generate(_ context.Context, req *types.GenerateRequest) (*types.GenerateResponse, error) {

	p := req.GetPayload()
	_, err := s.h.Write(p)
	if err != nil {
		return nil, err
	}
	sum := fmt.Sprintf("%x", s.h.Sum(p))

	id := uuid.NewV4()

	fmt.Println("sum", sum, "id", id)
	return &types.GenerateResponse{
		Md5: sum,
		Id:  id.Bytes(),
	}, nil
}

func (s *Server) handler(m *nats.Msg) {
	var req types.GenerateRequest
	proto.Unmarshal(m.Data, &req)
	resp, _ := s.Generate(context.TODO(), &req)
	b, _ := proto.Marshal(resp)
	s.nc.Publish(m.Reply, b)
}

func (s *Server) Run() {
	s.one.Do(func() {
		sub, _ := s.nc.QueueSubscribe(subjGenerate, "pubsub", s.handler)
		s.subs[subjGenerate] = sub
	})
}

func (s *Server) Close() {
	for _, sub := range s.subs {
		sub.Unsubscribe()
	}
}

func NewServer(nc *nats.Conn) *Server {
	return &Server{
		nc:   nc,
		one:  &sync.Once{},
		subs: make(map[string]*nats.Subscription),
		h:    md5.New(),
	}
}
