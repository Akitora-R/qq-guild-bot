package conn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/log"
	"github.com/tencent-connect/botgo/openapi"
	"qq-guild-bot/internal/conn/entity"
	"qq-guild-bot/internal/pkg/config"
	"strconv"
	"strings"
)

func GetSelf() dto.User {
	return *selfInfo
}

func PostMessage(content, msgId, channelId string, fileImage []byte) (entity.SendMessageResponse, error) {
	authToken := fmt.Sprintf("%s.%s", strconv.FormatUint(config.AppConf.AppID, 10), config.AppConf.AccessToken)
	req := resty.New().R().SetAuthScheme("Bot").SetAuthToken(authToken)
	if len(fileImage) > 0 {
		req.SetFileReader("file_image", "file_name.jpg", bytes.NewReader(fileImage))
	}
	r, err := req.SetMultipartFormData(map[string]string{
		"content": content,
		"msg_id":  msgId,
	}).SetPathParam("channel_id", channelId).Post(config.BaseApi + "/channels/{channel_id}/messages")
	var resp entity.SendMessageResponse
	if err != nil {
		return resp, err
	}
	respBody := r.Body()
	log.Info(string(respBody))
	err = json.Unmarshal(respBody, &resp)
	return resp, err
}

func CreateDirectMessage(sourceGuildId, recipientId string) (*dto.DirectMessage, error) {
	return botApi.CreateDirectMessage(botCtx, &dto.DirectMessageToCreate{
		SourceGuildID: sourceGuildId,
		RecipientID:   recipientId,
	})
}

func PostDirectMessage(guildId string, msg *dto.MessageToCreate) (*dto.Message, error) {
	return botApi.PostDirectMessage(botCtx, &dto.DirectMessage{
		GuildID: guildId,
	}, msg)
}

func GuildMembers(guildId, after string, limit int) ([]*dto.Member, error) {
	p := dto.GuildMembersPager{
		After: after,
		Limit: strconv.Itoa(limit),
	}
	return botApi.GuildMembers(botCtx, guildId, &p)
}

func DelMsg(channelId, messageId string, hideTip bool) error {
	var err error
	if hideTip {
		err = botApi.RetractMessage(botCtx, channelId, messageId, openapi.RetractMessageOptionHidetip)
	} else {
		err = botApi.RetractMessage(botCtx, channelId, messageId)
	}
	return err
}

func UpdateUser(u *entity.UserUpdate, guildId, userId string) error {
	if u.MuteSecond != nil {
		updateGuildMute := &dto.UpdateGuildMute{
			MuteSeconds: strconv.Itoa(*u.MuteSecond),
		}
		return botApi.MemberMute(botCtx, guildId, userId, updateGuildMute)
	}
	return nil
}

func GetGuildMember(guildId, userId string) (*dto.Member, error) {
	return botApi.GuildMember(botCtx, guildId, userId)
}

func GetRoles(guildId string) (*dto.GuildRoles, error) {
	return botApi.Roles(botCtx, guildId)
}

func AddMemberRole(guildId, roleId, userId string, b *dto.MemberAddRoleBody) error {
	return botApi.MemberAddRole(botCtx, guildId, dto.RoleID(roleId), userId, b)
}

func DeleteGuildMember(guildId, userId string, deleteHistoryMsgDay *int, addBlackList *bool) error {
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
	return botApi.DeleteGuildMember(botCtx, guildId, userId, args...)
}

func BanByBatch(guildId string, memberIdList []string) error {
	var failed []string
	log.Info("名单长度:", len(memberIdList))
	for _, mId := range memberIdList {
		log.Info("ban: ", mId)
		if err := botApi.DeleteGuildMember(botCtx, guildId, mId); err != nil {
			log.Error(err)
			failed = append(failed, mId)
		}
	}
	if len(failed) > 0 {
		return errors.New(strings.Join(failed, ","))
	}
	return nil
}
