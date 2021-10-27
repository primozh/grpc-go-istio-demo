package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pb "github.com/primozh/grpc-go/proto"
	"google.golang.org/grpc"
)

var serverName = os.Getenv("SERVER_NAME")

const (
	port = 8000
)

type helloServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *helloServer) Greeting(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: fmt.Sprintf("Hello %s from server %s", request.Name, serverName)}, nil
}

func (s *helloServer) Greetings(request *pb.HelloRequest, stream pb.HelloService_GreetingsServer) error {
	stream.Send(&pb.HelloResponse{Message: fmt.Sprintf("Hello %s from server %s", request.Name, serverName)})
	stream.Send(&pb.HelloResponse{Message: fmt.Sprintf("Hello %s from server %s", request.Name, serverName)})
	stream.Send(&pb.HelloResponse{Message: fmt.Sprintf("Hello %s from server %s", request.Name, serverName)})
	stream.Send(&pb.HelloResponse{Message: fmt.Sprintf("Hello %s from server %s", request.Name, serverName)})
	return nil
}

func (s *helloServer) GreetingsBi(stream pb.HelloService_GreetingsBiServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		stream.Send(&pb.HelloResponse{Message: fmt.Sprintf("Hello %s from server %s", in.Name, serverName)})
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))

	log.Printf("Starting server with params %s\n", serverName)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Listening...")

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	s := new(helloServer)
	pb.RegisterHelloServiceServer(grpcServer, s)
	grpcServer.Serve(lis)
}
