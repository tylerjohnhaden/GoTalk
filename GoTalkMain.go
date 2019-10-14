package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	pb "./proto"
)

const (
	rpcPort  = ":50051"
	httpPort = ":8080"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// SayHello implements helloworld.GreeterServer
func (s *server) SubmitJobRequest(stream *pb.JobRequest) {

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)
		... // look for notes to be sent to client
		for _, note := range s.routeNotes[key] {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	go log.Fatal(http.ListenAndServe(httpPort, nil))

	lis, err := net.Listen("tcp", rpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
