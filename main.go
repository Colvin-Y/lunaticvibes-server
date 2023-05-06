package main
import (
	"context"
	"log"
	"net"

	helloworldpb "github.com/Colvin-Y/lunaticvibes-server/proto/helloworld"
	"google.golang.org/grpc"
)

type server struct {
	helloworldpb.UnimplementedGreeterServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) SayHello(ctx context.Context, in *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	log.Println("recieve")
	return &helloworldpb.HelloReply{Message: in.Name + " world"}, nil
}

func main(){
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// 创建一个gRPC server对象
	s := grpc.NewServer()
	// 注册Greeter service到server
	helloworldpb.RegisterGreeterServer(s, &server{})
	// 8080端口启动gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	log.Fatalln(s.Serve(lis))
}