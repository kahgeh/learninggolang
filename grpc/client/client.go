package main

import (
	"context"
	pb "learngogrpc/hello"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:10000"
)

func main() {
	cnn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer cnn.Close()
	client := pb.NewHelloClient(cnn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayIt(ctx, &pb.SayItRequest{Language: "en"})
	if err != nil {
		log.Fatalf("could not SayIt: %v", err)
	}
	log.Println(r.Text)
}
