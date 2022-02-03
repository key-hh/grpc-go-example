package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/key-hh/grpc-go-example/grpc/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gRPCAddr = "0.0.0.0:8089"
	restAddr = "0.0.0.0:8090"
)

func main() {
	ctx := context.Background()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := hello.RegisterGreeterHandlerFromEndpoint(ctx, mux, gRPCAddr, opts)
	if err != nil {
		log.Fatalf("%v", err)
	}
	go http.ListenAndServe(restAddr, mux)

	log.Printf("gateway is ready for %s", restAddr)

	sigQuitCh := make(chan os.Signal)
	signal.Notify(sigQuitCh, syscall.SIGQUIT)
	<-sigQuitCh
}
