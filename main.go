package main

import (
	"log"
	"net"
	"os"

	helloworldpb "github.com/Colvin-Y/lunaticvibes-server/proto/helloworld"
	scorepb "github.com/Colvin-Y/lunaticvibes-server/proto/score"
	service "github.com/Colvin-Y/lunaticvibes-server/service"
	logger "github.com/Colvin-Y/lunaticvibes-server/utils/log"
	"google.golang.org/grpc"
)

func main() {
	// 初始化 logger
	logger, err := logger.NewLogger("/var/log/lunaticvibes/lc.log")
	if err != nil {
		os.Exit(1)
	}
	defer logger.Close()

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// 创建一个gRPC server对象
	s := grpc.NewServer()

	// 注册 service
	svc := &service.Server{}
	svc.Logger = logger
	helloworldpb.RegisterGreeterServer(s, svc)
	scorepb.RegisterInsertScoreServer(s, svc)

	// 8080端口启动gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	log.Fatalln(s.Serve(lis))
}
