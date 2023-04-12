package entity

import (
	"github.com/tencent-connect/botgo/dto"
)

type EventType string

const (
	Message           EventType = "MESSAGE"
	MessageDelete     EventType = "MESSAGE_DELETE"
	DirectMessage     EventType = "DIRECT_MESSAGE"
	GuildMemberAdd    EventType = "GUILD_MEMBER_ADD"
	GuildMemberUpdate EventType = "GUILD_MEMBER_UPDATE"
	GuildMemberRemove EventType = "GUILD_MEMBER_REMOVE"
)

type GuildEvent[T any] struct {
	Id        string    `json:"id,omitempty"`
	EventType EventType `json:"event_type"`
	Self      *dto.User `json:"self"`
	Data      *T        `json:"data"`
}

type MessageEventData struct {
	*dto.Message
}

type MessageDeleteEventData struct {
	Message *dto.Message `json:"message"`
	OpUser  *dto.User    `json:"op_user"`
}

type MemberEventData struct {
	*dto.Member
}

func NewMessageEvent(id string, self *dto.User, d *dto.Message) *GuildEvent[MessageEventData] {
	return &GuildEvent[MessageEventData]{
		Id:        id,
		EventType: Message,
		Self:      self,
		Data:      &MessageEventData{Message: d},
	}
}

func NewDirectMessageEvent(id string, self *dto.User, d *dto.Message) *GuildEvent[MessageEventData] {
	return &GuildEvent[MessageEventData]{
		Id:        id,
		EventType: DirectMessage,
		Self:      self,
		Data:      &MessageEventData{Message: d},
	}
}

func NewMessageDeleteEvent(id string, self, opUser *dto.User, d *dto.Message) *GuildEvent[MessageDeleteEventData] {
	return &GuildEvent[MessageDeleteEventData]{
		Id:        id,
		EventType: MessageDelete,
		Self:      self,
		Data:      &MessageDeleteEventData{Message: d, OpUser: opUser},
	}
}

func NewMemberAddEventData(id string, self *dto.User, member *dto.Member) *GuildEvent[MemberEventData] {
	return newMemberEventData(id, self, member, GuildMemberAdd)
}

func NewMemberUpdateEventData(id string, self *dto.User, member *dto.Member) *GuildEvent[MemberEventData] {
	return newMemberEventData(id, self, member, GuildMemberUpdate)
}

func NewMemberRemoveEventData(id string, self *dto.User, member *dto.Member) *GuildEvent[MemberEventData] {
	return newMemberEventData(id, self, member, GuildMemberRemove)
}

func newMemberEventData(id string, self *dto.User, member *dto.Member, eventType EventType) *GuildEvent[MemberEventData] {
	return &GuildEvent[MemberEventData]{
		Id:        id,
		EventType: eventType,
		Self:      self,
		Data:      &MemberEventData{Member: member},
	}
}
