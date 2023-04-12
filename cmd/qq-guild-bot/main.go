package main

import (
	"qq-guild-bot/internal/api"
	"qq-guild-bot/internal/conn"
)

func main() {
	conn.StartGuildEventListen()
	api.StartHttpAPI()
}
