package main

import (
	"context"
	"github.com/thomaszub/grpc-example/greet/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type server struct{}

func (s *server) GreetManyTimes(request *pb.GreetManyTimesRequest, timesServer pb.GreetService_GreetManyTimesServer) error {
	firstName := request.GetGreeting().FirstName
	lastName := request.GetGreeting().LastName
	result := "Hello " + firstName + " " + lastName + "!"
	for i := 0; i < 10; i++ {
		res := &pb.GreetManyTimesResponse{Result: result}
		timesServer.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (s *server) Greet(_ context.Context, request *pb.GreetRequest) (*pb.GreetResponse, error) {
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
