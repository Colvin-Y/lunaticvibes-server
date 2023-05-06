package service

import (
	"context"
	"log"

	helloworldpb "github.com/Colvin-Y/lunaticvibes-server/proto/helloworld"
)

func (s *Server) SayHello(ctx context.Context, in *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	log.Println("say hello")
	return &helloworldpb.HelloReply{Message: in.Name + " world"}, nil
}
