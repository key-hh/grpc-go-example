package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/key-hh/grpc-go-example/grpc/stream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpAddr = ":8091"
	grpcAddr = "localhost:8092"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := stream.NewStreamerClient(conn)

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		s, err := client.SendToStream(ctx)
		if err != nil {
			log.Printf("SendToStream err %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := 0; i < 5; i++ {
			err := s.Send(&stream.StreamRequest{
				Name: "client/SendToStream request" + strconv.Itoa(i),
			})
			if err != nil {
				log.Printf("SendToStream send err %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		reply, err := s.CloseAndRecv()
		if err != nil {
			log.Printf("SendToStream CloseAndRecv err %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("client/SendToStream -- receive: ", reply.Name, reply.Age)

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		s, err := client.ReceiveFromStream(ctx, &stream.StreamRequest{
			Name: "client/ReceiveFromStream request",
		})
		if err != nil {
			log.Printf("ReceiveFromStream err %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for {
			d, err := s.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("ReceiveFromStream err %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println("client/ReceiveFromStream -- receive: ", d.Name, d.Age)
		}
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		s, err := client.ChatInStream(ctx)
		if err != nil {
			log.Printf("ChatInStream err %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		waitCh := make(chan struct{})

		go func() {
			for {
				d, err := s.Recv()
				if err == io.EOF {
					close(waitCh)
					return
				}
				if err != nil {
					log.Printf("Chat ReceiveFromStream err %v", err)
					return
				}
				log.Println("client/ChatInStream -- receive: ", d.Name, d.Age)
			}
		}()

		for i := 0; i < 7; i++ {
			err := s.Send(&stream.StreamRequest{
				Name: "client/ChatInStream request" + strconv.Itoa(i),
			})
			if err != nil {
				log.Printf("ChatInStream send err %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		s.CloseSend()

		<-waitCh
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("server is ready for %s", httpAddr)
	http.ListenAndServe(httpAddr, nil)
}
