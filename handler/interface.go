package handler

import (
	"context"

	scorepb "github.com/Colvin-Y/lunaticvibes-server/proto/score"
	logger "github.com/Colvin-Y/lunaticvibes-server/utils/log"
)

type InsertScoreHandlerFunc func(ctx context.Context, handler *InsertScoreHandler) error

type InsertScoreHandler struct {
	Logger *logger.Logger
	Req    *scorepb.InsertScoreRequest
	Resp   *scorepb.InsertScoreReply
}

func (handler *InsertScoreHandler) Sync(ctx context.Context, funcs ...InsertScoreHandlerFunc) error {
	for _, funcItem := range funcs {
		err := funcItem(ctx, handler)
		if err != nil {
			handler.Logger.Error(err.Error())
			return err
		}
	}
	return nil
}
