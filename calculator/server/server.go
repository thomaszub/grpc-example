package main

import (
	"context"
	"github.com/thomaszub/grpc-example/calculator/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

type server struct{}

type invalidNumberError struct {
	value int64
}

func (e *invalidNumberError) Error() string {
	return "Invalid number error. Number " + strconv.Itoa(int(e.value)) + " is not greater than 2"
}

func (s *server) Sum(_ context.Context, request *pb.SumRequest) (*pb.SumResponse, error) {
	values := request.Values
	var result int64 = 0
	for _, v := range values {
		result += v
	}
	res := &pb.SumResponse{Result: result}
	return res, nil
}

func (s *server) PrimeNumbers(request *pb.PrimeNumberRequest, numbersServer pb.CalculatorService_PrimeNumbersServer) error {
	number := request.Value
	if number < 2 {
		return &invalidNumberError{value: number}
	}
	for number > 2 {
		var i int64
		for i = 2; i < 1000; i++ {
			if number%i == 0 {
				number = number / i
				numbersServer.Send(&pb.PrimeNumberResponse{Result: i})
				break
			}
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	serv := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(serv, &server{})

	if err := serv.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
