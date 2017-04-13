package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/zensh/grpc-bench/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "127.0.0.1:10000"
)

const total = 1000000
const cocurrency = 1000

func main() {
	tc, err := credentials.NewClientTLSFromFile("../secret/ca.crt", "127.0.0.1")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(tc))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	_, err = c.SayHello(context.Background(), &pb.HelloRequest{Name: "Ping"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	var w sync.WaitGroup
	w.Add(total)
	co := make(chan int, cocurrency)

	task := func(i int) {
		defer w.Done()
		_, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "Ping"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		<-co
		if (i % 1000) == 0 {
			fmt.Print(".")
		}
	}

	t := time.Now()
	for i := 0; i < total; i++ {
		co <- i
		go task(i)
	}
	log.Println("Wait")
	w.Wait()
	sec := time.Now().Sub(t) / 1e6
	log.Printf("\nFinished, cocurrency: %d, time: %d ms, %d ops", cocurrency, sec, (total / (sec / 1000)))
}
