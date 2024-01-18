package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	"qq-guild-bot/internal/pkg/stub"
)

type QQGuildServiceServerImpl struct {
	stub.UnimplementedQQGuildServiceServer
	ch chan *stub.GuildEventData
}

func NewQQGuildServer(ch chan *stub.GuildEventData) stub.QQGuildServiceServer {
	return &QQGuildServiceServerImpl{
		ch: ch,
	}
}

func (s *QQGuildServiceServerImpl) GetMessageStream(_ *empty.Empty, stream stub.QQGuildService_GetMessageStreamServer) error {
	for message := range s.ch {
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	return nil
}

func (s *QQGuildServiceServerImpl) PostEvent(_ context.Context, _ *stub.GuildEventData) (*emptypb.Empty, error) {
	panic("no need")
}

type QQGuildServiceClientImpl struct {
}

func (c *QQGuildServiceClientImpl) Post() {
}
