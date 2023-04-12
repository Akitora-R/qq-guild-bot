package entity

import "github.com/tencent-connect/botgo/dto"

type SendMessageResponse struct {
	ID           string        `json:"id"`
	ChannelID    string        `json:"channel_id"`
	GuildID      string        `json:"guild_id"`
	Content      string        `json:"content"`
	Timestamp    dto.Timestamp `json:"timestamp"`
	Author       dto.User      `json:"author"`
	D            bool          `json:"d"`
	Type         int           `json:"type"`
	Flags        int           `json:"flags"`
	SeqInChannel string        `json:"seq_in_channel"`
	Code         *int          `json:"code"`
	Message      *string       `json:"message"`
}
