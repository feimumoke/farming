package main

import (
	"context"
	"feimumoke/farming/v2/api/service"
	"feimumoke/farming/v2/framework"
	"feimumoke/farming/v2/service/identify"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", framework.NetGrpcAddress, "gRPC server endpoint")
)

func main() {
	go func() {
		runGrpc()
	}()
	if err := runGw(); err != nil {
		log.Fatal(err)
	}
}

func runGw() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := service.RegisterIdentifyServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(framework.NetGwAddress, mux)
}

func runGrpc() {
	lis, err := net.Listen(framework.NetProtocol, framework.NetGrpcAddress)
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	service.RegisterIdentifyServiceServer(s, identify.NewIdentifySvr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
