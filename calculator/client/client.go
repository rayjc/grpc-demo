package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/rayjc/grpc-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client connected.")

	url := fmt.Sprintf("localhost:%v", "50051")
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer conn.Close()

	client := calculatorpb.NewCaculatorServiceClient(conn)
	// fmt.Printf("Client created: %f", client)

	// runUnary(client)

	// runServerStream(client)

	runClientStream(client)
}

func runUnary(client calculatorpb.CaculatorServiceClient) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := r.Int() % 10
	values := make([]int32, length)
	for i := 0; i < length; i++ {
		values[i] = r.Int31() % 10
	}

	fmt.Printf("Calling a Unary RPC with {%v}...\n", values)
	req := &calculatorpb.SumRequest{
		Values: values,
	}

	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.Result)
}

func runServerStream(client calculatorpb.CaculatorServiceClient) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Int63() % 100
	fmt.Printf("Calling a Server stream RPC with {%v}...\n", num)
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: num,
	}

	stream, err := client.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Prime Number Decomposition RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unexpected error occured during server streaming...")
		}
		// parse response
		log.Printf("Response from Prime Number Decomposition RPC: %v", res.GetPrimeFactor())
	}
}

func runClientStream(client calculatorpb.CaculatorServiceClient) {
	size := 5
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	requests := []*calculatorpb.AverageRequest{}
	for i := 0; i < size; i++ {
		requests = append(requests, &calculatorpb.AverageRequest{
			Value: int64(r.Int31()),
		})
	}
	fmt.Printf("Calling a Client stream RPC with {%v}...\n", requests)

	stream, err := client.Average(context.Background())
	if err != nil {
		log.Fatalf("Error calling Average: %v\n", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving Average response: %v\n", err)
	}
	fmt.Printf("Average response: %v\n", res.Result)
}
