package service

import (
	helloworldpb "github.com/Colvin-Y/lunaticvibes-server/proto/helloworld"
	scorepb "github.com/Colvin-Y/lunaticvibes-server/proto/score"
)

type Server struct {
	helloworldpb.UnimplementedGreeterServer
	scorepb.UnimplementedInsertScoreServer
}

func NewServer() *Server {
	return &Server{}
}
