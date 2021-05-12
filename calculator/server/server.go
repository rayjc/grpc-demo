package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/rayjc/grpc-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function is called with %v\n", req)
	values := req.GetValues()
	var result int32 = 0
	for _, v := range values {
		result += v
	}

	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CaculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition is called with %v\n", req)
	num := req.GetNumber()
	// prime number decomposition
	var divisor int64 = 2
	for num > 1 {
		if num%divisor == 0 {
			// send prime number
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			num /= divisor
		} else {
			divisor++
		}
	}
	return nil
}

func main() {
	fmt.Println("Connected.")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCaculatorServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
