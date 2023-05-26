package main

import (
	"context"
	"flag"
	"fmt"
	proxy2 "github.com/Vingurzhou/zwz-proxy/proxy"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"

	gw "github.com/Vingurzhou/zwz-proxy/proto" // Update
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterHelloWorldHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	addr := fmt.Sprintf(":%d", 8081)
	log.Printf("server listening at [::]%v", addr)
	proxy := proxy2.NewProxy()
	http.Handle("/", proxy)
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	for _, address := range addresses {
		ipNet, ok := address.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				log.Printf("server running at %s", ipNet.IP.String())
			}
		}
	}
	return http.ListenAndServe(addr, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
