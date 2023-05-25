package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Vingurzhou/zwz-proxy/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9090, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedHelloWorldServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) Echo(ctx context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	log.Printf("Received: %v", in.Value)
	return &pb.StringMessage{Value: "Hello " + in.Value}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloWorldServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
