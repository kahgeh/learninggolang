package main

import (
	"context"
	"flag"
	"fmt"
	pb "learngogrpc/hello"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
)

type myHelloServer struct {
	defaultLanguage string
	lookup          map[string]string
}

func (s *myHelloServer) SayIt(context context.Context, request *pb.SayItRequest) (*pb.SayItResponse, error) {
	text := s.lookup[request.Language]
	return &pb.SayItResponse{Text: text, Language: request.Language}, nil
}

func newServer() *myHelloServer {
	server := &myHelloServer{
		lookup: make(map[string]string),
	}
	server.lookup["en"] = "g'day"
	server.lookup["cn"] = "你好"
	return server
}

func main() {
	flag.Parse()
	fmt.Printf("port = %d \n", *port)
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterHelloServer(grpcServer, newServer())
	grpcServer.Serve(listener)
}
