package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/key-hh/grpc-go-example/grpc/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpAddr = ":8088"
	grpcAddr = "localhost:8089"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := hello.NewGreeterClient(conn)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := new(hello.HelloRequest)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp, err := client.SayHello(ctx, req)
		if err != nil {
			log.Printf("SayHello error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("SayHello resp: %v", resp)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("SayHello error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	log.Printf("server is ready for %s", httpAddr)
	http.ListenAndServe(httpAddr, nil)
}
