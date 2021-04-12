package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rayjc/grpc-demo/greet/config"
	"github.com/rayjc/grpc-demo/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client connected.")

	url := fmt.Sprintf("localhost:%v", config.Port)
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)
	// fmt.Printf("Client created: %f", client)

	runUnary(client)
}

func runUnary(client greetpb.GreetServiceClient) {
	greeting := &greetpb.Greeting{
		FirstName: "Bruce",
		LastName:  "Wayne",
	}
	fmt.Printf("Calling a Unary RPC with {%v}...\n", greeting)
	req := &greetpb.GreetRequest{
		Greeting: greeting,
	}

	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}
