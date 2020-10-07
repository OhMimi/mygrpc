package main

import (
	"context"
	"log"

	pb "mygrpc/echospec"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:9999", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("連線失敗 : %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoServerClient(conn)

	r, err := c.Echo(context.Background(), &pb.EchoRequest{
		Message: "Hi Hi Hi",
	})
	if err != nil {
		log.Fatalf("無法執行 Plus 函式：%v", err)
	}

	log.Printf("回傳結果：%s , 時間:%d", r.Message, r.Unixtime)
}
