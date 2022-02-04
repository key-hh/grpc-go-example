package handler

import (
	"io"
	"log"

	"github.com/key-hh/grpc-go-example/grpc/stream"
)

type StreamerServer struct {
	stream.UnimplementedStreamerServer
}

func (StreamerServer) SendToStream(s stream.Streamer_SendToStreamServer) error {
	recvCount := 0
	for {
		d, err := s.Recv()
		if err == io.EOF {
			return s.SendAndClose(&stream.StreamResponse{
				Name: "SendToStream response",
				Age:  int32(recvCount),
			})
		}
		if err != nil {
			return err
		}

		recvCount++

		log.Println("SendToStream -- receive: ", d.Name, recvCount)
	}
}

func (StreamerServer) ReceiveFromStream(r *stream.StreamRequest, s stream.Streamer_ReceiveFromStreamServer) error {
	log.Println("ReceiveFromStream -- receive: ", r.Name)

	for i := 0; i < 10; i++ {
		s.Send(&stream.StreamResponse{
			Name: "ReceiveFromStream response",
			Age:  int32(i),
		})
	}
	return nil
}

func (StreamerServer) ChatInStream(s stream.Streamer_ChatInStreamServer) error {
	recvCount := 0
	for {
		d, err := s.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		recvCount++
		log.Println("ChatInStream -- receive: ", d.Name, recvCount)

		if err := s.Send(&stream.StreamResponse{
			Name: "ChatInStream response",
			Age:  int32(recvCount),
		}); err != nil {
			return err
		}
	}
}
