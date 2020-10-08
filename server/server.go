package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	pb "mygrpc/echospec"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 準備幾個fake users
var users = map[int32]pb.UserResponse{
	1: {
		Name: "HuanNa",
		Age:  20,
	},
	2: {
		Name: "Kevin",
		Age:  30,
	},
}

type Server struct{}

func (s *Server) Echo(c context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Printf("receive client request, client send:%s\n", req.Message)
	return &pb.EchoResponse{
		Message:  req.Message,
		Unixtime: time.Now().Unix(),
	}, nil
}

func (s *Server) GetUserInfo(c context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	// 查找map有沒有該user, 有就回覆, 否則就回錯誤
	if user, ok := users[req.GetID()]; ok {
		return &user, nil
	}

	log.Printf("req : %v\n", req)
	return nil, errors.New("user not found")
}

func (s *Server) GetUserStreamInfo(stream pb.UserService_GetUserStreamInfoServer) error {
	// 開一個for 一直收或是發, 直到我們自己想離開為止
	for {
		// 透過Recv(), 從stream收取cleint打來的資料
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatal(err)
			return err
		}

		if user, ok := users[req.GetID()]; ok {
			// server主動送回給client
			err = stream.Send(&user)
			if err != nil {
				log.Fatal(err)
				return err
			}
		} else {
			log.Printf("req : %v\n", req)
			return fmt.Errorf("user not found: %d\n", req.GetID())
		}
	}
}

func main() {

	// 註冊端口來提供gRPC服務
	apiListener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
		return
	}

	// 註冊grpc
	es := &Server{}

	// 建構一個gRPC服務端實例
	grpc := grpc.NewServer()

	// 註冊服務
	pb.RegisterEchoServerServer(grpc, es)
	pb.RegisterUserServiceServer(grpc, es)

	reflection.Register(grpc)

	// grpc server服務啟用
	if err := grpc.Serve(apiListener); err != nil {
		log.Fatal(" grpc.Serve Error: ", err)
		return
	}

}
