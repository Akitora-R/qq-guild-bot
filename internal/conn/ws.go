package conn

import (
	"context"
	"errors"
	"qq-guild-bot/internal/conn/entity"
	"qq-guild-bot/internal/pkg/config"
	"qq-guild-bot/internal/pkg/util/http"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/log"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

var botApi openapi.OpenAPI
var botCtx context.Context
var selfInfo *dto.User

const bufferSize = 100

/*
https://github.com/tencent-connect/botgo/tree/master/examples
*/

func StartGuildEventListen() {
	conf := config.AppConf
	go func(conf config.Config) {
		ch := make(chan any, bufferSize)
		startGuildBotWS(ch, conf)
		for data := range ch {
			for _, rep := range conf.ServerConfigs {
				go func(d any, server config.ServerConfig) {
					_, _, _ = http.NewPostRequest(server.Url).SetBodyJson(d).Exec()
				}(data, rep)
			}
		}
	}(conf)
}

func startGuildBotWS(ch chan any, conf config.Config) {
	go func(cha chan any) {
		t := token.BotToken(conf.AppID, conf.AccessToken)
		if conf.Sandbox {
			log.Info("沙箱环境")
			botApi = botgo.NewSandboxOpenAPI(t).WithTimeout(3 * time.Second)
		} else {
			log.Info("正式环境")
			botApi = botgo.NewOpenAPI(t).WithTimeout(3 * time.Second)
		}
		botCtx = context.Background()
		ws, err := botApi.WS(botCtx, nil, "")
		if err != nil {
			panic(err)
		}
		selfInfo, err = botApi.Me(botCtx)
		if err != nil {
			log.Warn("请求自身信息失败", err)
		}
		log.Info("self: ", selfInfo)
		intent := websocket.RegisterHandlers(
			messageEventHandler(cha),
			directMessageEventHandler(cha),
			messageDeleteEventHandler(cha),
			memberEventHandler(cha),
			messageReactionEventHandler(cha),
			interactionEventHandler(cha),
		)
		err = botgo.NewSessionManager().Start(ws, t, &intent)
		if err != nil {
			panic(err)
		}
	}(ch)
}

func messageEventHandler(ch chan any) event.MessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageData) error {
		log.Info("消息id: ", data.ID)
		ch <- entity.NewMessageEvent(event.Id, selfInfo, (*dto.Message)(data))
		return nil
	}
}

func directMessageEventHandler(ch chan any) event.DirectMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSDirectMessageData) error {
		log.Info("私信id: ", data.ID)
		ch <- entity.NewDirectMessageEvent(event.Id, selfInfo, (*dto.Message)(data))
		return nil
	}
}

func messageDeleteEventHandler(ch chan any) event.MessageDeleteEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageDeleteData) error {
		ch <- entity.NewMessageDeleteEvent(event.Id, selfInfo, &data.OpUser, &data.Message)
		return nil
	}
}

func memberEventHandler(ch chan any) event.GuildMemberEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildMemberData) error {
		switch event.Type {
		case "GUILD_MEMBER_ADD":
			ch <- entity.NewMemberAddEventData(event.Id, selfInfo, (*dto.Member)(data))
		case "GUILD_MEMBER_UPDATE":
			ch <- entity.NewMemberUpdateEventData(event.Id, selfInfo, (*dto.Member)(data))
		case "GUILD_MEMBER_REMOVE":
			ch <- entity.NewMemberRemoveEventData(event.Id, selfInfo, (*dto.Member)(data))
		default:
			return errors.New("未知类型：" + string(event.Type))
		}
		return nil
	}
}

func messageReactionEventHandler(ch chan any) event.MessageReactionEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageReactionData) error {
		var t = entity.EventType(event.Type)
		ch <- entity.NewMessageReactionEventData(event.Id, selfInfo, t, (*dto.MessageReaction)(data))
		return nil
	}
}

func interactionEventHandler(ch chan any) event.InteractionEventHandler {
	return func(event *dto.WSPayload, data *dto.WSInteractionData) error {
		log.Info(string(event.RawMessage))
		ch <- entity.NewInteractionEventData(event.Id, selfInfo, (*dto.Interaction)(data))
		return nil
	}
}
