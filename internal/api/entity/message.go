package entity

import (
	"encoding/base64"
	"github.com/tencent-connect/botgo/dto/message"
	"strconv"
	"strings"
	"time"
)

type MessageType string

const (
	Text        MessageType = "text"
	MentionUser MessageType = "at"
	Emoji       MessageType = "emoji"
)

type Message struct {
	EventId     string        `json:"event_id,omitempty"`
	MsgId       string        `json:"msg_id,omitempty"`
	Messages    []MessageData `json:"messages"`
	Attachments []Attachment  `json:"attachments"`
}

func (m Message) ToContent() string {
	var content []string
	for _, data := range m.Messages {
		content = append(content, data.toContent())
	}
	return strings.Join(content, " ")
}

type MessageData struct {
	Type MessageType `json:"type,omitempty"`
	Data string      `json:"data,omitempty"`
}

func (d MessageData) toContent() string {
	switch d.Type {
	case Text:
		return d.Data
	case MentionUser:
		return message.MentionUser(d.Data)
	case Emoji:
		i, err := strconv.Atoi(d.Data)
		if err != nil {
			return ""
		}
		return message.Emoji(i)
	}
	return ""
}

type Attachment struct {
	Url    string `json:"url"`
	Base64 string `json:"base64"`
}

func (a Attachment) ToBytes() []byte {
	if b, err := base64.StdEncoding.DecodeString(a.Base64); err == nil {
		return b
	}
	return nil
}

type MsgResp struct {
	Id              string    `json:"id"`
	ChannelId       string    `json:"channel_id"`
	GuildId         string    `json:"guild_id"`
	Content         string    `json:"content"`
	Timestamp       time.Time `json:"timestamp"`
	Tts             bool      `json:"tts"`
	MentionEveryone bool      `json:"mention_everyone"`
	Author          struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		Bot      bool   `json:"bot"`
	} `json:"author"`
	Pinned       bool   `json:"pinned"`
	Type         int    `json:"type"`
	Flags        int    `json:"flags"`
	SeqInChannel string `json:"seq_in_channel"`
}
