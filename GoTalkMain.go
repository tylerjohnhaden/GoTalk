package main

import (
	pb "./gotalk"
	"context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	rpcPort = ":50051"
	//httpPort = ":8080"
	rpcAddress = "localhost:50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

//var (
//	gotalkServer *Server
//)

//func handler(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
//}

// SayHello implements helloworld.GreeterServer
func (s *server) SubmitJobRequest(ctx context.Context, jobRequest *pb.JobRequest) (*pb.JobResponse, error) {
	log.Printf("Processing job request")

	newData := map[string]string{
		"favoriteIceCream": "vanilla",
	}

	for k, v := range jobRequest.GetJobData() {
		newData[k] = v
	}

	return &pb.JobResponse{
		JobData: newData,
	}, nil
}

func serverStuff() {
	lis, err := net.Listen("tcp", rpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGoTalkServer(s, &server{})
	log.Print("Starting server")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func clientStuff() {
	conn, err := grpc.Dial(rpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect client: %v", err)
	}
	defer conn.Close()

	c := pb.NewGoTalkClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Print("Starting client")

	r, err := c.SubmitJobRequest(ctx, &pb.JobRequest{
		JobId: rand.Uint64(),
		JobNode: &pb.JobNode{
			JobTitle: pb.JobNode_A,
			JobNodes: []*pb.JobNode{
				{JobTitle: pb.JobNode_B,},
				{JobTitle: pb.JobNode_C,},
				{JobTitle: pb.JobNode_D,},
				{JobTitle: pb.JobNode_E,},
				{JobTitle: pb.JobNode_F,},
			},
		},
		JobData: map[string]string{
			"firstName": "tyler",
			"lastName": "haden",
		},
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Recieved Response Data: %s", r.JobData)
}

func main() {
	//http.HandleFunc("/", handler)
	//go log.Fatal(http.ListenAndServe(httpPort, nil))

	go serverStuff()

	time.Sleep(2 * time.Second)

	clientStuff()
}
