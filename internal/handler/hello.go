package handler

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"

	"github.com/key-hh/grpc-go-example/grpc/hello"
	"google.golang.org/protobuf/types/known/anypb"
	wpb "google.golang.org/protobuf/types/known/wrapperspb"
)

type GreeterServer struct {
	hello.GreeterServer
}

func (*GreeterServer) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	log.Printf("request: %v", req)

	details := new(anypb.Any)
	err := anypb.MarshalFrom(details, wpb.String("hello, world"), proto.MarshalOptions{})
	if err != nil {
		return nil, err
	}

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
		Details: details,
		Projects: map[string]int32{
			"key1": 11,
			"key2": 12,
		},
	}, nil
}
