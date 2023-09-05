package conn

import (
	"errors"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"log/slog"
	"qq-guild-bot/internal/conn/entity"
	"qq-guild-bot/internal/pkg/config"
	"qq-guild-bot/internal/pkg/util/http"
	"time"
)

const bufferSize = 100

var bots = map[string]*Bot{}

// StartGuildEventListen https://github.com/tencent-connect/botgo/tree/master/examples
func StartGuildEventListen() {
	conf := config.AppConf
	go func(conf config.Config) {
		ch := make(chan any, bufferSize)
		for _, botConfig := range conf.Bot {
			bot := NewBot(botConfig, ch)
			go func() {
				if err := bot.Start(); err != nil {
					slog.Error("start bot failed", err, bot)
				}
			}()
			for i := 0; i < 10; i++ {
				if bot.GetSelf() != nil {
					bots[bot.selfInfo.ID] = bot
					slog.Info("bot login", "id", bot.selfInfo.ID)
					break
				}
				time.Sleep(time.Second)
			}
		}
		for data := range ch {
			for _, rep := range conf.Server {
				go func(d any, server config.ServerConfig) {
					_, _, _ = http.NewPostRequest(server.Url).SetBodyJson(d).Exec()
				}(data, rep)
			}
		}
	}(conf)
}

func (b *Bot) messageEventHandler() event.MessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageData) error {
		slog.Info("收到消息", "bot", b.selfInfo.ID, "id", data.ID, "content", data.Content)
		b.ch <- entity.NewMessageEvent(event.Id, b.selfInfo, (*dto.Message)(data))
		return nil
	}
}

func (b *Bot) directMessageEventHandler() event.DirectMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSDirectMessageData) error {
		slog.Info("私信id: ", data.ID)
		b.ch <- entity.NewDirectMessageEvent(event.Id, b.selfInfo, (*dto.Message)(data))
		return nil
	}
}

func (b *Bot) messageDeleteEventHandler() event.MessageDeleteEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageDeleteData) error {
		b.ch <- entity.NewMessageDeleteEvent(event.Id, b.selfInfo, &data.OpUser, &data.Message)
		return nil
	}
}

func (b *Bot) memberEventHandler() event.GuildMemberEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildMemberData) error {
		switch event.Type {
		case "GUILD_MEMBER_ADD":
			b.ch <- entity.NewMemberAddEventData(event.Id, b.selfInfo, (*dto.Member)(data))
		case "GUILD_MEMBER_UPDATE":
			b.ch <- entity.NewMemberUpdateEventData(event.Id, b.selfInfo, (*dto.Member)(data))
		case "GUILD_MEMBER_REMOVE":
			b.ch <- entity.NewMemberRemoveEventData(event.Id, b.selfInfo, (*dto.Member)(data))
		default:
			return errors.New("未知类型：" + string(event.Type))
		}
		return nil
	}
}

func (b *Bot) messageReactionEventHandler() event.MessageReactionEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageReactionData) error {
		var t = entity.EventType(event.Type)
		b.ch <- entity.NewMessageReactionEventData(event.Id, b.selfInfo, t, (*dto.MessageReaction)(data))
		return nil
	}
}

func (b *Bot) interactionEventHandler() event.InteractionEventHandler {
	return func(event *dto.WSPayload, data *dto.WSInteractionData) error {
		slog.Info(string(event.RawMessage))
		b.ch <- entity.NewInteractionEventData(event.Id, b.selfInfo, (*dto.Interaction)(data))
		return nil
	}
}
