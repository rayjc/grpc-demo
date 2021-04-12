package main

import (
	"context"
	"fmt"
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

	runUnary(client)
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
