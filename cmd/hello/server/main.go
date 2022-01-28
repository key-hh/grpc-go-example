package main

import (
	"fmt"
	"log"
	"net"

	"github.com/key-hh/grpc-go-example/grpc/hello"
	"github.com/key-hh/grpc-go-example/internal/handler"
	"google.golang.org/grpc"
)

const (
	host = "0.0.0.0"
	port = 8089
)

func main() {
	addr := fmt.Sprintf("%s:%d", host, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	hello.RegisterGreeterServer(srv, &handler.GreeterServer{})

	log.Printf("server is ready for %s", addr)

	if err := srv.Serve(ln); err != nil {
		if err == grpc.ErrServerStopped {
			log.Printf("ListenAndServe: %v", err)
		} else {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}
}
