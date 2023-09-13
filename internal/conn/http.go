package conn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
	"log/slog"
	apiEntity "qq-guild-bot/internal/api/entity"
	"qq-guild-bot/internal/conn/entity"
	"strconv"
	"strings"
)

func (b *Bot) GetSelf() *dto.User {
	return b.selfInfo
}

func (b *Bot) PostMessage(channelId string, m *apiEntity.Message, fileImage []byte) (entity.SendMessageResponse, error) {
	authToken := fmt.Sprintf("%s.%s", strconv.FormatUint(b.config.AppID, 10), b.config.AccessToken)
	req := resty.New().R().SetAuthScheme("Bot").SetAuthToken(authToken)
	if len(fileImage) > 0 {
		req.SetFileReader("file_image", "file_name.jpg", bytes.NewReader(fileImage))
	}
	r, err := req.SetMultipartFormData(map[string]string{
		"content":  m.ToContent(),
		"msg_id":   m.MsgId,
		"event_id": m.EventId,
	}).SetPathParam("channel_id", channelId).Post(b.config.Endpoint + "/channels/{channel_id}/messages")
	var resp entity.SendMessageResponse
	if err != nil {
		return resp, err
	}
	respBody := r.Body()
	err = json.Unmarshal(respBody, &resp)
	return resp, err
}

func (b *Bot) CreateDirectMessage(sourceGuildId, recipientId string) (*dto.DirectMessage, error) {
	return b.api.CreateDirectMessage(b.ctx, &dto.DirectMessageToCreate{
		SourceGuildID: sourceGuildId,
		RecipientID:   recipientId,
	})
}

func (b *Bot) PostDirectMessage(guildId string, msg *dto.MessageToCreate) (*dto.Message, error) {
	return b.api.PostDirectMessage(b.ctx, &dto.DirectMessage{
		GuildID: guildId,
	}, msg)
}

func (b *Bot) GuildMembers(guildId, after string, limit int) ([]*dto.Member, error) {
	p := dto.GuildMembersPager{
		After: after,
		Limit: strconv.Itoa(limit),
	}
	return b.api.GuildMembers(b.ctx, guildId, &p)
}

func (b *Bot) DelMsg(channelId, messageId string, hideTip bool) error {
	var err error
	if hideTip {
		err = b.api.RetractMessage(b.ctx, channelId, messageId, openapi.RetractMessageOptionHidetip)
	} else {
		err = b.api.RetractMessage(b.ctx, channelId, messageId)
	}
	return err
}

func (b *Bot) UpdateUser(u *entity.UserUpdate, guildId, userId string) error {
	if u.MuteSecond != nil {
		updateGuildMute := &dto.UpdateGuildMute{
			MuteSeconds: strconv.Itoa(*u.MuteSecond),
		}
		return b.api.MemberMute(b.ctx, guildId, userId, updateGuildMute)
	}
	return nil
}

func (b *Bot) GetGuildMember(guildId, userId string) (*dto.Member, error) {
	return b.api.GuildMember(b.ctx, guildId, userId)
}

func (b *Bot) GetRoles(guildId string) (*dto.GuildRoles, error) {
	return b.api.Roles(b.ctx, guildId)
}

func (b *Bot) AddMemberRole(guildId, roleId, userId string, body *dto.MemberAddRoleBody) error {
	return b.api.MemberAddRole(b.ctx, guildId, dto.RoleID(roleId), userId, body)
}

func (b *Bot) DeleteMemberRole(guildId, roleId, userId string, body *dto.MemberAddRoleBody) error {
	return b.api.MemberDeleteRole(b.ctx, guildId, dto.RoleID(roleId), userId, body)
}

func (b *Bot) DeleteGuildMember(guildId, userId string, deleteHistoryMsgDay *int, addBlackList *bool) error {
	var args []dto.MemberDeleteOption
	if deleteHistoryMsgDay != nil {
		args = append(args, func(opts *dto.MemberDeleteOpts) {
			opts.DeleteHistoryMsgDays = *deleteHistoryMsgDay
		})
	}
	if addBlackList != nil {
		args = append(args, func(opts *dto.MemberDeleteOpts) {
			opts.AddBlackList = *addBlackList
		})
	}
	return b.api.DeleteGuildMember(b.ctx, guildId, userId, args...)
}

func (b *Bot) BanByBatch(guildId string, memberIdList []string) error {
	var failed []string
	slog.Info("", "名单长度", len(memberIdList))
	for _, mId := range memberIdList {
		slog.Info("", "ban", mId)
		if err := b.api.DeleteGuildMember(b.ctx, guildId, mId); err != nil {
			slog.Error("", "err", err)
			failed = append(failed, mId)
		}
	}
	if len(failed) > 0 {
		return errors.New(strings.Join(failed, ","))
	}
	return nil
}
