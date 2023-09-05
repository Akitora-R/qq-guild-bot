package conn

import (
	"context"
	"errors"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"log/slog"
	"qq-guild-bot/internal/pkg/config"
	"time"
)

type Bot struct {
	config   config.BotConfig
	ctx      context.Context
	api      openapi.OpenAPI
	ch       chan any
	selfInfo *dto.User
}

func NewBot(config config.BotConfig, ch chan any) *Bot {
	t := token.BotToken(config.AppID, config.AccessToken)
	ctx := context.Background()
	var api openapi.OpenAPI
	if config.Sandbox {
		api = botgo.NewSandboxOpenAPI(t).WithTimeout(3 * time.Second)
	} else {
		api = botgo.NewOpenAPI(t).WithTimeout(3 * time.Second)
	}
	b := Bot{
		config: config,
		ctx:    ctx,
		ch:     ch,
		api:    api,
	}
	return &b
}

func (b *Bot) Start() error {
	ws, err := b.api.WS(b.ctx, nil, "")
	if err != nil {
		return err
	}
	selfInfo, err := b.api.Me(b.ctx)
	if err != nil {
		return err
	}
	b.selfInfo = selfInfo
	slog.Info("bot started", "self", b.selfInfo)
	handler := event.NewHandler(
		b.messageEventHandler(),
		b.directMessageEventHandler(),
		b.messageDeleteEventHandler(),
		b.memberEventHandler(),
		b.messageReactionEventHandler(),
		b.interactionEventHandler(),
	)
	t := token.BotToken(b.config.AppID, b.config.AccessToken)
	return botgo.NewSessionManager().Start(ws, t, handler)
}

func GetBotInstance(id string) (*Bot, error) {
	if id == "" {
		if len(bots) <= 0 {
			return nil, errors.New("no available bot instance")
		}
		for _, bot := range bots {
			return bot, nil
		}
	}
	bot := bots[id]
	if bot == nil {
		return nil, errors.New("invalid self id: " + id)
	}
	return bot, nil
}

func GetBotInstances() map[string]*Bot {
	return bots
}
