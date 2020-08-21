package main

import (
	"context"
	"github.com/thomaszub/grpc-example/calculator/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
}

func (s server) Sum(_ context.Context, request *pb.SumRequest) (*pb.SumResponse, error) {
	values := request.Values
	var result int64 = 0
	for _, v := range values {
		result += v
	}
	res := &pb.SumResponse{Result: result}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	serv := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(serv, server{})

	if err := serv.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
