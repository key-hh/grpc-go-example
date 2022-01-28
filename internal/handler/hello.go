package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/key-hh/grpc-go-example/grpc/hello"
	"google.golang.org/protobuf/types/known/anypb"
)

type GreeterServer struct {
	hello.GreeterServer
}

func (*GreeterServer) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	log.Printf("request: %v", req)

	return &hello.HelloReply{
		Message: fmt.Sprintf("%s result message", req.Name),
		Results: []*hello.HelloReply_Result{
			{
				Name: req.Name,
				Test: req.Age + 1,
			},
			{
				Name: req.Name,
				Test: req.Age + 2,
			},
		},
		Details: &anypb.Any{TypeUrl: "", Value: []byte("test detail with any")},
		Projects: map[string]int32{
			"key1": 11,
			"key2": 12,
		},
	}, nil
}
