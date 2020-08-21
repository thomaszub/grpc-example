package main

import (
	"context"
	"fmt"
	"github.com/thomaszub/grpc-example/greet/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not build connection: %v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Error in closing connection: %v", err)
		}
	}()

	cl := pb.NewGreetServiceClient(conn)
	doUnary(cl)
	doServerStreaming(cl)
	doClientStreaming(cl)
}

func doUnary(cl pb.GreetServiceClient) {
	gr := &pb.Greeting{
		FirstName: "Thomas",
		LastName:  "Zub",
	}
	req := &pb.GreetRequest{Greeting: gr}
	res, err := cl.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error fetching response: %v", err)
	}

	log.Printf("Got response: %q", res.Result)
}

func doServerStreaming(cl pb.GreetServiceClient) {
	gr := &pb.Greeting{
		FirstName: "Thomas",
		LastName:  "Zub",
	}
	req := &pb.GreetManyTimesRequest{Greeting: gr}
	gCl, err := cl.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling RPC: %v", err)
	}
	for {
		msg, err := gCl.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching response: %v", err)
		}
		log.Printf("Got response: %q", msg.Result)
	}
}

func doClientStreaming(cl pb.GreetServiceClient) {
	greetings := []*pb.Greeting{
		{
			FirstName: "Thomas",
			LastName:  "Zub",
		},
		{
			FirstName: "Brian",
			LastName:  "Adams",
		},
		{
			FirstName: "Luke",
			LastName:  "Skywalker",
		},
	}

	gCl, err := cl.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error calling RPC: %v", err)
	}
	for _, gr := range greetings {
		req := &pb.LongGreetRequest{Greeting: gr}
		err := gCl.Send(req)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	res, err := gCl.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}
	fmt.Printf("Got result: %q", res.Result)
}
