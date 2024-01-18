package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	grpcLib "google.golang.org/grpc"
	"net"
	"qq-guild-bot/internal/api/grpc"
	"qq-guild-bot/internal/api/http"
	"qq-guild-bot/internal/conn"
	"qq-guild-bot/internal/conn/entity"
	"qq-guild-bot/internal/pkg"
	"qq-guild-bot/internal/pkg/config"
	"qq-guild-bot/internal/pkg/stub"
	"time"
)

func main() {
	httpCh := conn.GetEventChan()
	grpcCh := make(chan *stub.GuildEventData)
	go handleGuildEvent(httpCh, grpcCh)
	s := grpcLib.NewServer()
	if err := pkg.RegisterConsul(s); err != nil {
		panic(err)
	}
	gs := grpc.NewQQGuildServer(grpcCh)
	stub.RegisterQQGuildServiceServer(s, gs)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.AppConf.GrpcPort))
	if err != nil {
		panic(err)
	}
	go http.StartHttpAPI()
	if err = s.Serve(lis); err != nil {
		panic(err)
	}
}

func handleGuildEvent(ch chan entity.GuildEvent, grpcCh chan *stub.GuildEventData) {
	for data := range ch {
		for _, rep := range config.AppConf.Server {
			go func(d any, server config.ServerConfig) {
				_, _ = resty.New().R().SetBody(d).Post(server.Url)
			}(data, rep)
		}
		select {
		case grpcCh <- convert(data):
		case <-time.After(100 * time.Millisecond):
		}
	}
}

func convert(event entity.GuildEvent) *stub.GuildEventData {
	data := stub.GuildEventData{
		Id:        event.GetId(),
		EventType: stub.EventType(stub.EventType_value[string(event.GetEventType())]),
		Self:      grpc.ConvertUser(event.GetSelf()),
		Data:      nil,
	}
	switch eventType := event.GetEventType(); eventType {
	case entity.DirectMessage:
		fallthrough
	case entity.Message:
		data.Data = grpc.ConvertMessageEventData(event.GetData().(*entity.MessageEventData))
	case entity.MessageDelete:
		data.Data = grpc.ConvertMessageDeleteEventData(event.GetData().(*entity.MessageDeleteEventData))
	case entity.GuildMemberAdd:
		fallthrough
	case entity.GuildMemberUpdate:
		fallthrough
	case entity.GuildMemberRemove:
		data.Data = grpc.ConvertMemberEventData(event.GetData().(*entity.MemberEventData))
	case entity.MessageReactionAdd:
		fallthrough
	case entity.MessageReactionRemove:
		data.Data = grpc.ConvertMessageReactionEventData(event.GetData().(*entity.MessageReactionEventData))
	case entity.InteractionCreate:
		data.Data = grpc.ConvertInteractionEventData(event.GetData().(*entity.InteractionEventData))
	}
	return &data
}
