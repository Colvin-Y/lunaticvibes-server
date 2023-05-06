package service

import (
	"context"

	"github.com/Colvin-Y/lunaticvibes-server/handler"
	scorepb "github.com/Colvin-Y/lunaticvibes-server/proto/score"
)

func precheck(ctx context.Context, handler *handler.InsertScoreHandler) error {
	err := handler.Req.ValidateAll()
	if err != nil {
		handler.Logger.Error(err.Error())
		return err
	}

	return nil
}

func (s *Server) InsertScore(ctx context.Context, in *scorepb.InsertScoreRequest) (*scorepb.InsertScoreReply, error) {
	hd := &handler.InsertScoreHandler{
		Logger: s.Logger,
		Req:    in,
		Resp:   &scorepb.InsertScoreReply{},
	}

	funcs := []handler.InsertScoreHandlerFunc{
		precheck,
	}

	err := hd.Sync(ctx, funcs...)
	var respMsg string
	if err != nil {
		respMsg = err.Error()
	} else {
		respMsg = in.Data.String()
	}
	return &scorepb.InsertScoreReply{Message: respMsg}, nil
}
