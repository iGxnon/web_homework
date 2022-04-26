//client（调用rpc的一方）
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework16/level-1/server/proto"
	"log"
)

const (
	address = "localhost:50051"
)

func main() {
	//建立链接
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewBiliClient(conn)

	for {
		//这段不重要
		fmt.Println("input username&password:")
		iptName := ""
		_, _ = fmt.Scanln(&iptName)
		iptPassword := ""
		_, _ = fmt.Scanln(&iptPassword)

		loginResp, _ := c.Login(context.Background(), &proto.LoginReq{
			UserName: iptName,
			PassWord: iptPassword,
		})

		if loginResp.Status == proto.Status_OK {
			fmt.Println("success")
			break
		}
		fmt.Println("retry")
	}
}
