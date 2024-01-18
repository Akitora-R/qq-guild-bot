package grpc

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/tencent-connect/botgo/dto"
	"qq-guild-bot/internal/conn/entity"
	"qq-guild-bot/internal/pkg/stub"
)

func ConvertMessageEventData(data *entity.MessageEventData) *stub.GuildEventData_MessageEventData {
	return &stub.GuildEventData_MessageEventData{MessageEventData: convertMessage(data.Message)}
}

func ConvertMessageDeleteEventData(data *entity.MessageDeleteEventData) *stub.GuildEventData_MessageDeleteEventData {
	return &stub.GuildEventData_MessageDeleteEventData{
		MessageDeleteEventData: &stub.MessageDeleteEventData{
			Message: convertMessage(data.Message),
			OpUser:  ConvertUser(data.OpUser),
		},
	}
}

func ConvertMemberEventData(data *entity.MemberEventData) *stub.GuildEventData_MemberEventData {
	return &stub.GuildEventData_MemberEventData{
		MemberEventData: convertMember(data.Member),
	}
}

func ConvertMessageReactionEventData(data *entity.MessageReactionEventData) *stub.GuildEventData_MessageReactionEventData {
	return &stub.GuildEventData_MessageReactionEventData{
		MessageReactionEventData: convertMessageReaction(data),
	}
}

func ConvertInteractionEventData(data *entity.InteractionEventData) *stub.GuildEventData_InteractionEventData {
	interactionData := stub.InteractionData{
		Name:     data.Data.Name,
		Type:     stub.InteractionDataType(data.Data.Type),
		Resolved: data.Data.Resolved,
	}
	return &stub.GuildEventData_InteractionEventData{
		InteractionEventData: &stub.Interaction{
			Id:            data.ID,
			ApplicationId: data.ApplicationID,
			Type:          stub.InteractionType(data.Type),
			Data:          &interactionData,
			GuildId:       data.GuildID,
			ChannelId:     data.ChannelID,
			Version:       data.Version,
		},
	}
}

func convertMessageReaction(data *entity.MessageReactionEventData) *stub.MessageReaction {
	reactionTarget := stub.ReactionTarget{
		Id:   data.Target.ID,
		Type: stub.ReactionTargetType(data.Target.Type),
	}
	emoji := stub.Emoji{
		Id:   data.Emoji.ID,
		Type: int32(data.Emoji.Type),
	}
	return &stub.MessageReaction{
		UserId:    data.UserID,
		ChannelId: data.ChannelID,
		GuildId:   data.GuildID,
		Target:    &reactionTarget,
		Emoji:     &emoji,
	}
}

func convertMessage(data *dto.Message) *stub.Message {
	return &stub.Message{
		Id:               data.ID,
		ChannelId:        data.ChannelID,
		GuildId:          data.GuildID,
		Content:          data.Content,
		Timestamp:        convertTimestamp(data.Timestamp),
		EditedTimestamp:  convertTimestamp(data.EditedTimestamp),
		MentionEveryone:  data.MentionEveryone,
		Author:           ConvertUser(data.Author),
		Member:           convertMember(data.Member),
		Attachments:      convertAttachments(data.Attachments),
		Embeds:           convertEmbeds(data.Embeds),
		Mentions:         convertUsers(data.Mentions),
		Ark:              convertArk(data.Ark),
		DirectMessage:    data.DirectMessage,
		SeqInChannel:     data.SeqInChannel,
		MessageReference: convertMessageReference(data.MessageReference),
		SrcGuildId:       data.SrcGuildID,
	}
}

func convertMessageReference(r *dto.MessageReference) *stub.MessageReference {
	if r == nil {
		return nil
	}
	return &stub.MessageReference{
		MessageId:             r.MessageID,
		IgnoreGetMessageError: r.IgnoreGetMessageError,
	}
}

func convertArk(ark *dto.Ark) *stub.Ark {
	if ark == nil {
		return nil
	}
	return &stub.Ark{
		TemplateId: int32(ark.TemplateID),
		Kv:         convertArkKVs(ark.KV),
	}
}

func convertArkKVs(kvs []*dto.ArkKV) []*stub.ArkKV {
	if len(kvs) <= 0 {
		return nil
	}
	var r []*stub.ArkKV
	for _, kv := range kvs {
		r = append(r, &stub.ArkKV{
			Key:   kv.Key,
			Value: kv.Value,
			Obj:   convertArkObjs(kv.Obj),
		})
	}
	return r
}

func convertArkObjs(objs []*dto.ArkObj) []*stub.ArkObj {
	if len(objs) <= 0 {
		return nil
	}
	var r []*stub.ArkObj
	for _, kv := range objs {
		r = append(r, &stub.ArkObj{
			ObjKv: convertArkObjKVs(kv.ObjKV),
		})
	}
	return r
}

func convertArkObjKVs(objKVs []*dto.ArkObjKV) []*stub.ArkObjKV {
	var r []*stub.ArkObjKV
	for _, kv := range objKVs {
		r = append(r, &stub.ArkObjKV{
			Key:   kv.Key,
			Value: kv.Value,
		})
	}
	return r
}

func convertMember(member *dto.Member) *stub.Member {
	return &stub.Member{
		GuildId:  member.GuildID,
		JoinedAt: convertTimestamp(member.JoinedAt),
		Nick:     member.Nick,
		User:     ConvertUser(member.User),
		Roles:    member.Roles,
		OpUserId: member.OpUserID,
	}
}

func convertUsers(users []*dto.User) []*stub.User {
	var r []*stub.User
	for _, u := range users {
		r = append(r, ConvertUser(u))
	}
	return r
}

func ConvertUser(user *dto.User) *stub.User {
	if user == nil {
		return nil
	}
	return &stub.User{
		Id:               user.ID,
		Username:         user.Username,
		Avatar:           user.Avatar,
		Bot:              user.Bot,
		UnionOpenid:      user.UnionOpenID,
		UnionUserAccount: user.UnionUserAccount,
	}
}

func convertEmbeds(embeds []*dto.Embed) []*stub.Embed {
	var r []*stub.Embed
	for _, embed := range embeds {
		var fields []*stub.EmbedField
		for _, field := range fields {
			fields = append(fields, &stub.EmbedField{
				Name:  field.Name,
				Value: field.Value,
			})
		}
		r = append(r, &stub.Embed{
			Title:       embed.Title,
			Description: embed.Description,
			Prompt:      embed.Prompt,
			Thumbnail:   &stub.MessageEmbedThumbnail{Url: embed.Thumbnail.URL},
			Fields:      fields,
		})
	}
	return nil
}

func convertAttachments(as []*dto.MessageAttachment) []*stub.MessageAttachment {
	var gas []*stub.MessageAttachment
	for _, a := range as {
		gas = append(gas, &stub.MessageAttachment{Url: a.URL})
	}
	return gas
}

// convertTimestamp converts time.Time to protobuf Timestamp
func convertTimestamp(time dto.Timestamp) *timestamp.Timestamp {
	t, _ := time.Time()
	if t.IsZero() {
		return nil
	}
	return &timestamp.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}
