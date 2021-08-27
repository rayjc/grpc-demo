package main

import (
	"context"
	"fmt"
	"io"
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

func (*server) Average(stream calculatorpb.CaculatorService_AverageServer) error {
	fmt.Println("Average is called.")
	var total int64 = 0
	var count int64 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			var result int64 = 0
			if count != 0 {
				result = total / count
			}
			stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: result,
			})
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		total += req.GetValue()
		count++
	}
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
