package main

import (
	"context"
	"github.com/thomaszub/grpc-example/calculator/pb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not build connection: %v", err)
	}
	cl := pb.NewCalculatorServiceClient(conn)
	//doUnary(cl)
	doServerStreaming(cl)
}

func doUnary(cl pb.CalculatorServiceClient) {
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

func doServerStreaming(cl pb.CalculatorServiceClient) {
	var value int64 = 231
	req := &pb.PrimeNumberRequest{Value: value}
	pCl, err := cl.PrimeNumbers(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling RPC: %v", err)
	}

	for {
		res, err := pCl.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error fetching response: %v", err)
		}
		log.Printf("Got response: %d", res.Result)
	}
}
