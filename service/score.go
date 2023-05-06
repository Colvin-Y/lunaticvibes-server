package service

import (
	"context"
	"log"

	scorepb "github.com/Colvin-Y/lunaticvibes-server/proto/score"
)

func (s *Server) InsertScore(ctx context.Context, in *scorepb.InsertScoreRequest) (*scorepb.InsertScoreReply, error) {
	log.Println("insert score")
	err := in.ValidateAll()
	if err != nil {
		log.Println(err)
		return &scorepb.InsertScoreReply{Message: err.Error()}, nil
	}
	return &scorepb.InsertScoreReply{Message: in.Data.String()}, nil
}
