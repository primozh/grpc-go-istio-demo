package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "github.com/primozh/grpc-go/proto"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ticker(client pb.HelloServiceClient) {
	tick := time.Tick(5 * time.Second)
	for range tick {
		response, err := client.Greeting(context.Background(), &pb.HelloRequest{Name: "Primoz"})
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Unary response from server: %s", response.Message)
		}

		stream, err := client.Greetings(context.Background(), &pb.HelloRequest{Name: "Primoz"})

		if err != nil {
			log.Fatal(err)
		}
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error %v", err)
			}

			log.Printf("Server stream response: %s", response.Message)
		}

		names := []*pb.HelloRequest{
			{Name: "Primoz"},
			{Name: "Marko"},
			{Name: "Matej"},
		}
		bistream, err := client.GreetingsBi(context.Background())

		waitc := make(chan struct{})
		go func() {
			for {
				in, err := bistream.Recv()
				if err == io.EOF {
					close(waitc)
					return
				}
				if err != nil {
					log.Fatalf("Failed to receive a note: %v", err)
				}
				log.Printf("Bistream server response: %s", in.Message)
			}
		}()

		for _, name := range names {
			if err := bistream.Send(name); err != nil {
				log.Fatalf("Failed to send a name: %v", err)
			}
		}
		bistream.CloseSend()
		<-waitc
	}
}

func main() {
	serverAddress := getEnv("SERVER_ADDRESS", "localhost")
	serverPort := getEnv("SERVER_PORT", "8080")
	connString := fmt.Sprintf("%s:%s", serverAddress, serverPort)

	log.Printf("Connecting to server on %s", connString)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(connString, opts...)
	if err != nil {
		panic(err)
	}
	client := pb.NewHelloServiceClient(conn)
	ticker(client)
}
