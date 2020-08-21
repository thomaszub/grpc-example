package main

import (
	"context"
	"github.com/thomaszub/grpc-example/greet/pb"
	"google.golang.org/grpc"
	"io"
	"log"
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
