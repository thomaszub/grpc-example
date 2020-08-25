package main

import (
	"context"
	"errors"
	"github.com/thomaszub/grpc-example/calculator/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"math"
	"net"
	"strconv"
)

type server struct{}

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
		return errors.New("Invalid number error. Number " + strconv.Itoa(int(number)) + " is not greater than 2")
	}
	for number > 2 {
		var i int64
		for i = 2; i < 1000; i++ {
			if number%i == 0 {
				number = number / i
				res := &pb.PrimeNumberResponse{Result: i}
				err := numbersServer.Send(res)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

func (s *server) ComputeAverage(averageServer pb.CalculatorService_ComputeAverageServer) error {
	var count int64 = 0
	var sum int64 = 0
	for {
		req, err := averageServer.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sum += req.Value
		count++
	}
	if count == 0 {
		// No values from stream, returning 0
		res := &pb.ComputeAverageResponse{Result: 0.0}
		err := averageServer.SendAndClose(res)
		if err != nil {
			return err
		}
	}
	result := float64(sum) / float64(count)
	res := &pb.ComputeAverageResponse{Result: result}
	err := averageServer.SendAndClose(res)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) FindMaximum(maximumServer pb.CalculatorService_FindMaximumServer) error {
	var max int64 = math.MinInt64
	for {
		res, err := maximumServer.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		val := res.Value
		if val > max {
			max = val
			res := &pb.FindMaximumResponse{CurrentMaximum: max}
			if err := maximumServer.Send(res); err != nil {
				return err
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
