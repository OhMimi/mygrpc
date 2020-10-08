package main

import (
	"context"
	"log"

	pb "mygrpc/echospec"

	"google.golang.org/grpc"
)

func main() {

	// 透過Dial()負責跟gRPC服務端建立起連線
	conn, err := grpc.Dial("localhost:9999", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("連線失敗 : %v", err)
	}
	defer conn.Close()

	EchoReq(conn)

	UserReq(conn)
}

func EchoReq(conn *grpc.ClientConn) {
	// 注入連線, 返回NewEchoServerClient對象
	c := pb.NewEchoServerClient(conn)

	r, err := c.Echo(context.Background(), &pb.EchoRequest{
		Message: "Hi Hi Hi",
	})
	if err != nil {
		log.Fatalf("無法執行 Plus 函式：%v", err)
	}

	log.Printf("回傳結果：%s , 時間:%d", r.Message, r.Unixtime)
}

func UserReq(conn *grpc.ClientConn) {
	// 注入連線, 返回NewEchoServerClient對象
	c := pb.NewUserServiceClient(conn)

	r, err := c.GetUserInfo(context.Background(), &pb.UserRequest{
		ID: 1,
	})
	if err != nil {
		log.Fatalf("無法執行 Plus 函式：%v", err)
	}

	log.Printf("回傳使用者名稱：%s , 年紀:%d", r.Name, r.Age)
}
