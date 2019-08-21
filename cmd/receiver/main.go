package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/sah4ez/grpc2nats/pkg/rpc"
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

	s := rpc.NewServer(nc)
	go s.Run()
	defer s.Close()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(
		shutdown,
		syscall.SIGHUP,
		syscall.SIGTERM,
	)

	<-shutdown

}
