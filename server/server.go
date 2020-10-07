package main

import (
	"context"
	"log"
	pb "mygrpc/echospec"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type EchoServer struct{}

func (e *EchoServer) Echo(c context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Printf("receive client request, client send:%s\n", req.Message)
	return &pb.EchoResponse{
		Message:  req.Message,
		Unixtime: time.Now().Unix(),
	}, nil
}

func main() {

	apiListener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
		return
	}

	// 註冊grpc
	es := &EchoServer{}

	grpc := grpc.NewServer()
	pb.RegisterEchoServerServer(grpc, es)

	reflection.Register(grpc)

	if err := grpc.Serve(apiListener); err != nil {
		log.Fatal(" grpc.Serve Error: ", err)
		return
	}

}
