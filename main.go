package main

import (
	"log"
	"net"

	helloworldpb "github.com/Colvin-Y/lunaticvibes-server/proto/helloworld"
	scorepb "github.com/Colvin-Y/lunaticvibes-server/proto/score"
	service "github.com/Colvin-Y/lunaticvibes-server/service"
	"google.golang.org/grpc"
)

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// 创建一个gRPC server对象
	s := grpc.NewServer()
	// 注册Greeter service到server
	helloworldpb.RegisterGreeterServer(s, &service.Server{})
	scorepb.RegisterInsertScoreServer(s, &service.Server{})
	// 8080端口启动gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	log.Fatalln(s.Serve(lis))
}
