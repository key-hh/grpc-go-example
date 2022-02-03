package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/key-hh/grpc-go-example/grpc/hello"
	"github.com/key-hh/grpc-go-example/internal/handler"
	"google.golang.org/grpc"
)

const (
	gRPCAddr = "0.0.0.0:8089"
)

func main() {
	ln, err := net.Listen("tcp", gRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	hello.RegisterGreeterServer(srv, &handler.GreeterServer{})

	go func() {
		if err := srv.Serve(ln); err != nil {
			if err == grpc.ErrServerStopped {
				log.Printf("ListenAndServe: %v", err)
			} else {
				log.Fatalf("ListenAndServe: %v", err)
			}
		}
	}()

	log.Printf("server is ready for %s", gRPCAddr)

	sigQuitCh := make(chan os.Signal)
	signal.Notify(sigQuitCh, syscall.SIGQUIT)
	<-sigQuitCh

	srv.GracefulStop()
}
