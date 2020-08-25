package main

import (
	"context"
	"github.com/thomaszub/grpc-example/greet/pb"
	"google.golang.org/grpc"
	"io"
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
		err := timesServer.Send(res)
		if err != nil {
			return err
		}
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

func reduce(values []string) string {
	result := ""
	if len(values) == 0 {
		return result
	}
	for i := 0; i < len(values)-1; i++ {
		result += values[i] + ", "
	}
	result += values[len(values)-1]
	return result
}

func (s *server) LongGreet(greetServer pb.GreetService_LongGreetServer) error {
	var greets []string
	for {
		req, err := greetServer.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		greets = append(greets, req.Greeting.FirstName)
	}
	res := &pb.LongGreetResponse{Result: "Hello " + reduce(greets) + "!"}
	err := greetServer.SendAndClose(res)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) GreetEveryone(everyoneServer pb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := everyoneServer.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		result := "Hello " + req.Greeting.FirstName + "!"
		res := &pb.GreetEveryoneResponse{Result: result}
		if err := everyoneServer.Send(res); err != nil {
			return err
		}
	}
	return nil
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
