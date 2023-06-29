package service

import (
	helloworldpb "github.com/Colvin-Y/lunaticvibes-server/proto/helloworld"
	scorepb "github.com/Colvin-Y/lunaticvibes-server/proto/score"
	logger "github.com/Colvin-Y/lunaticvibes-server/utils/log"
)

type Server struct {
	helloworldpb.UnimplementedGreeterServer
	scorepb.UnimplementedScoreServer
	Logger *logger.Logger
}

func NewServer() *Server {
	return &Server{}
}
