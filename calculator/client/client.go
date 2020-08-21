package main

import (
	"context"
	"github.com/thomaszub/grpc-example/calculator/pb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not build connection: %v", err)
	}
	cl := pb.NewCalculatorServiceClient(conn)
	values := []int64{
		1, 2, 3, 4,
	}
	req := &pb.SumRequest{Values: values}
	val, err := cl.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error fetching response: %v", err)
	}
	log.Printf("Got response: %d", val.Result)
}
