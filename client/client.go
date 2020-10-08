package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	UserStreamReq(conn)
}

func EchoReq(conn *grpc.ClientConn) {
	// 注入連線, 返回NewEchoServerClient對象
	client := pb.NewEchoServerClient(conn)

	res, err := client.Echo(context.Background(), &pb.EchoRequest{
		Message: "Hi Hi Hi",
	})
	if err != nil {
		log.Fatalf("無法執行 Plus 函式：%v", err)
	}

	log.Printf("回傳結果：%s , 時間:%d", res.Message, res.Unixtime)
}

func UserReq(conn *grpc.ClientConn) {
	// 注入連線, 返回NewEchoServerClient對象
	client := pb.NewUserServiceClient(conn)

	res, err := client.GetUserInfo(context.Background(), &pb.UserRequest{
		ID: 1,
	})
	if err != nil {
		log.Fatalf("無法執行 Plus 函式：%v", err)
	}

	log.Printf("回傳使用者名稱：%s , 年紀:%d", res.Name, res.Age)
}

func UserStreamReq(conn *grpc.ClientConn) {
	// 返回一個userServiceClient實例, 它實現了UserServiceClient接口
	client := pb.NewUserServiceClient(conn)

	stream, err := client.GetUserStreamInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var userID int32
	for userID = 1; userID < 4; userID++ {
		// 發送多筆
		stream.Send(&pb.UserRequest{
			ID: userID,
		})
		fmt.Println("send: ", userID)
	}
	fmt.Println("send finish")
	time.Sleep(1 * time.Second)
	fmt.Println("start receive")
	for {
		reply, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("reply : %v\n", reply)
		}
	}
}
