package grpc

import (
	"github.com/golang/protobuf/ptypes/empty"
	"qq-guild-bot/internal/pkg/stub"
)

type QQGuildServerImpl struct {
	stub.UnimplementedQQGuildServer
	ch chan *stub.GuildEventData
}

func NewQQGuildServer(ch chan *stub.GuildEventData) stub.QQGuildServer {
	return &QQGuildServerImpl{
		ch: ch,
	}
}

func (s *QQGuildServerImpl) GetMessage(_ *empty.Empty, stream stub.QQGuild_GetMessageServer) error {
	for message := range s.ch {
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	return nil
}
