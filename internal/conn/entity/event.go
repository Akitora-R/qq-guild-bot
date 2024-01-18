package entity

import (
	"github.com/tencent-connect/botgo/dto"
)

type EventType string

const (
	Message               EventType = "MESSAGE"
	MessageDelete         EventType = "MESSAGE_DELETE"
	DirectMessage         EventType = "DIRECT_MESSAGE"
	GuildMemberAdd        EventType = "GUILD_MEMBER_ADD"
	GuildMemberUpdate     EventType = "GUILD_MEMBER_UPDATE"
	GuildMemberRemove     EventType = "GUILD_MEMBER_REMOVE"
	MessageReactionAdd    EventType = "MESSAGE_REACTION_ADD"
	MessageReactionRemove EventType = "MESSAGE_REACTION_REMOVE"
	InteractionCreate     EventType = "INTERACTION_CREATE"
)

type GuildEventMessageDataConstrain interface {
	MessageEventData |
		MessageDeleteEventData |
		MemberEventData |
		MessageReactionEventData |
		InteractionEventData
}

type GuildEvent interface {
	GetId() string
	GetEventType() EventType
	GetSelf() *dto.User
	GetData() any
}

type GuildEventData[T GuildEventMessageDataConstrain] struct {
	Id        string    `json:"id,omitempty"`
	EventType EventType `json:"event_type"`
	Self      *dto.User `json:"self"`
	Data      *T        `json:"data"`
}

func (g *GuildEventData[T]) GetId() string {
	return g.Id
}

func (g *GuildEventData[T]) GetEventType() EventType {
	return g.EventType
}

func (g *GuildEventData[T]) GetSelf() *dto.User {
	return g.Self
}

func (g *GuildEventData[T]) GetData() any {
	return g.Data
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

type MessageReactionEventData struct {
	*dto.MessageReaction
}

type InteractionEventData struct {
	*dto.Interaction
}

func NewMessageEvent(id string, self *dto.User, d *dto.Message) *GuildEventData[MessageEventData] {
	return &GuildEventData[MessageEventData]{
		Id:        id,
		EventType: Message,
		Self:      self,
		Data:      &MessageEventData{Message: d},
	}
}

func NewDirectMessageEvent(id string, self *dto.User, d *dto.Message) *GuildEventData[MessageEventData] {
	return &GuildEventData[MessageEventData]{
		Id:        id,
		EventType: DirectMessage,
		Self:      self,
		Data:      &MessageEventData{Message: d},
	}
}

func NewMessageDeleteEvent(id string, self, opUser *dto.User, d *dto.Message) *GuildEventData[MessageDeleteEventData] {
	return &GuildEventData[MessageDeleteEventData]{
		Id:        id,
		EventType: MessageDelete,
		Self:      self,
		Data:      &MessageDeleteEventData{Message: d, OpUser: opUser},
	}
}

func NewMemberAddEventData(id string, self *dto.User, member *dto.Member) *GuildEventData[MemberEventData] {
	return newMemberEventData(id, self, member, GuildMemberAdd)
}

func NewMemberUpdateEventData(id string, self *dto.User, member *dto.Member) *GuildEventData[MemberEventData] {
	return newMemberEventData(id, self, member, GuildMemberUpdate)
}

func NewMemberRemoveEventData(id string, self *dto.User, member *dto.Member) *GuildEventData[MemberEventData] {
	return newMemberEventData(id, self, member, GuildMemberRemove)
}

func newMemberEventData(id string, self *dto.User, member *dto.Member, eventType EventType) *GuildEventData[MemberEventData] {
	return &GuildEventData[MemberEventData]{
		Id:        id,
		EventType: eventType,
		Self:      self,
		Data:      &MemberEventData{Member: member},
	}
}

func NewMessageReactionEventData(id string, self *dto.User, eventType EventType, data *dto.MessageReaction) *GuildEventData[MessageReactionEventData] {
	return &GuildEventData[MessageReactionEventData]{
		Id:        id,
		EventType: eventType,
		Self:      self,
		Data:      &MessageReactionEventData{MessageReaction: data},
	}
}

func NewInteractionEventData(id string, self *dto.User, data *dto.Interaction) *GuildEventData[InteractionEventData] {
	return &GuildEventData[InteractionEventData]{
		Id:        id,
		EventType: InteractionCreate,
		Self:      self,
		Data:      &InteractionEventData{Interaction: data},
	}
}
