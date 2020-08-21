package main

import (
	"context"
	"github.com/thomaszub/grpc-example/greet/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s server) Greet(ctx context.Context, request *pb.GreetRequest) (*pb.GreetResponse, error) {
	firstName := request.GetGreeting().FirstName
	lastName := request.GetGreeting().LastName
	result := "Hello " + firstName + " " + lastName + "!"
	res := &pb.GreetResponse{Result: result}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
