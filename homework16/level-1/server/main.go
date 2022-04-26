//server（被调用rpc的一方）
package main

import (
	"context"
	"google.golang.org/grpc"
	proto2 "homework16/level-1/server/proto"
	"log"
	"net"
)

const (
	port = ":50051"
)

func main() {
	// 监听端口
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer() //获取新服务示例
	proto2.RegisterBiliServer(s, &server{})

	// 开始处理
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	proto2.UnimplementedBiliServer // 用于实现proto包里BiliServer接口
}

func (s *server) Login(ctx context.Context, req *proto2.LoginReq) (*proto2.LoginResp, error) {
	resp := &proto2.LoginResp{}
	log.Println("recv:", req.UserName, req.PassWord)
	if req.PassWord != GetPassWord(req.UserName) {
		resp.Status = proto2.Status_BAD
		return resp, nil
	}
	resp.Status = proto2.Status_OK
	return resp, nil
}

func GetPassWord(userName string) (password string) {
	return userName + "123456"
}
