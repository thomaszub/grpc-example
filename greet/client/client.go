package main

import (
	"context"
	"github.com/thomaszub/grpc-example/greet/pb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not build connection: %v", err)
	}
	defer conn.Close()

	cl := pb.NewGreetServiceClient(conn)
	doUnary(cl)
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
