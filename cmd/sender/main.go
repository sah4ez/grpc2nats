package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"bitbucket.org/ronte/msg-gokit-packages/types/uuid"
	nats "github.com/nats-io/nats.go"
	"github.com/sah4ez/grpc2nats/pkg/rpc"
	"github.com/sah4ez/grpc2nats/pkg/types"
)

func main() {
	var natsURL = nats.DefaultURL
	if len(os.Args) == 2 {
		natsURL = os.Args[1]
	}
	// Connect to the NATS server.
	nc, err := nats.Connect(natsURL, nats.Timeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	message := strings.Join(os.Args[len(os.Args)-1:], " ")

	c := rpc.NewClient(nc)
	resp, err := c.Generate(context.TODO(), &types.GenerateRequest{
		Payload: []byte(message),
	})
	if err != nil {
		fmt.Println("generate got error:", err.Error())
		os.Exit(2)
	}
	fmt.Println("md5", resp.GetMd5(), "id", uuid.FromBytesOrNil(resp.GetId()))
}
