package main

import (
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
	log.Printf("Created client: %f", cl)

}
