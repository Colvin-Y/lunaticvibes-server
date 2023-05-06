package service

import (
	"context"
	"fmt"

	constant "github.com/Colvin-Y/lunaticvibes-server/const"
	"github.com/Colvin-Y/lunaticvibes-server/handler"
	scorepb "github.com/Colvin-Y/lunaticvibes-server/proto/score"
	"github.com/Colvin-Y/lunaticvibes-server/utils/common"
	"go.mongodb.org/mongo-driver/bson"
)

func precheck(ctx context.Context, handler *handler.InsertScoreHandler) error {
	err := handler.Req.ValidateAll()
	if err != nil {
		handler.Logger.Error(err.Error())
		return err
	}

	return nil
}

func dataSave(ctx context.Context, handler *handler.InsertScoreHandler) error {
	handler.Logger.Info(fmt.Sprintf("start to insert data[%v] to mongo", handler.Req.Data.String()))
	collection, client, err := common.ConnectMongoTable(constant.DATABASE_NAME, constant.SCORE_TABLE)
	if err != nil {
		handler.Logger.Error(err.Error())
		return err
	}
	defer client.Disconnect(context.Background())

	data, err := bson.MarshalWithRegistry(bson.DefaultRegistry, handler.Req.Data)
	if err != nil {
		handler.Logger.Error(err.Error())
		return err
	}

	_, err = collection.InsertOne(context.Background(), data)
	if err != nil {
		handler.Logger.Error(err.Error())
		return err
	}
	handler.Logger.Info("insert score success")

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
		dataSave,
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
