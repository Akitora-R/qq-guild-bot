package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tencent-connect/botgo/dto"
	"html/template"
	"io"
	apiEntity "qq-guild-bot/internal/api/entity"
	"qq-guild-bot/internal/conn"
	connEntity "qq-guild-bot/internal/conn/entity"
	"qq-guild-bot/internal/embeded"
	myLogger "qq-guild-bot/internal/pkg/log"
	"qq-guild-bot/internal/pkg/util"
	"strconv"
)

func StartHttpAPI() {
	engine := gin.New()
	engine.Use(gin.Recovery(), myLogger.GinLoggerMiddleware())
	temp := template.Must(template.New("").Delims("[[", "]]").ParseFS(embeded.WebFiles, "web/template/*.html"))
	engine.SetHTMLTemplate(temp)
	engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{
			"title": "Main website",
		})
	})
	group := engine.Group("/api")
	group.GET("/me", func(c *gin.Context) {
		botInstance, err := getBotInstance(c)
		if err != nil {
			handleErr(c, err, nil)
			return
		}
		c.JSON(200, botInstance.GetSelf())
	})
	guildApi := group.Group("/guild/:guildId")
	guildApi.GET("/member", handle(getPagedMember))
	guildApi.GET("/member/:userId", handle(getMemberDetail))
	guildApi.PATCH("/member/:userId", handle(updateMember))
	guildApi.PUT("/member/:userId/roles/:roleId", handle(updateMemberRole))
	guildApi.DELETE("/member/:userId/roles/:roleId", handle(deleteMemberRole))
	guildApi.DELETE("/member/:userId", handle(deleteMember))
	guildApi.PUT("/direct-message/:userId", handle(createDirectMessage))
	guildApi.POST("/direct-message", handle(postDirectMessage))
	guildApi.GET("/roles", handle(getRoles))

	channelApi := group.Group("/channel/:channelId")
	channelApi.POST("/message", handle(sendMsg))
	channelApi.DELETE("/message/:messageId", handle(delMsg))
	if err := engine.Run("127.0.0.1:6800"); err != nil {
		panic(err)
	}
}

func getBotInstance(c *gin.Context) (*conn.Bot, error) {
	selfAppId, err := strconv.ParseUint(c.GetHeader("self"), 10, 64)
	if err != nil {
		if len(conn.Bots) <= 0 {
			return nil, errors.New("no available bot instance")
		}
		for _, bot := range conn.Bots {
			return bot, nil
		}
	}
	return conn.Bots[selfAppId], nil
}

func handle(handler func(c *gin.Context, bot *conn.Bot)) func(c *gin.Context) {
	return func(c *gin.Context) {
		botInstance, err := getBotInstance(c)
		if err != nil {
			handleErr(c, err, nil)
			return
		}
		handler(c, botInstance)
	}
}

func handleErr(c *gin.Context, err error, data any) {
	if err != nil {
		s := err.Error()
		c.JSON(500, apiEntity.NewErrResp[any](nil, &s))
	} else {
		c.JSON(200, apiEntity.NewOkResp[any](data, nil))
	}
}

func getRoles(c *gin.Context, bot *conn.Bot) {
	roles, err := bot.GetRoles(c.Param("guildId"))
	handleErr(c, err, roles)
}

func getMemberDetail(c *gin.Context, bot *conn.Bot) {
	guildId := c.Param("guildId")
	userId := c.Param("userId")
	member, err := bot.GetGuildMember(guildId, userId)
	handleErr(c, err, member)
}

func updateMemberRole(c *gin.Context, bot *conn.Bot) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	var b dto.MemberAddRoleBody
	if err == nil {
		_ = json.Unmarshal(bodyBytes, &b)
	}
	err = bot.AddMemberRole(c.Param("guildId"), c.Param("roleId"), c.Param("userId"), &b)
	handleErr(c, err, nil)
}

func deleteMemberRole(c *gin.Context, bot *conn.Bot) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	var b dto.MemberAddRoleBody
	if err == nil {
		_ = json.Unmarshal(bodyBytes, &b)
	}
	err = bot.DeleteMemberRole(c.Param("guildId"), c.Param("roleId"), c.Param("userId"), &b)
	handleErr(c, err, nil)
}

func updateMember(c *gin.Context, bot *conn.Bot) {
	u := util.MustParseReader[connEntity.UserUpdate](c.Request.Body)
	guildId := c.Param("guildId")
	userId := c.Param("userId")
	handleErr(c, bot.UpdateUser(u, guildId, userId), nil)
}

func sendMsg(c *gin.Context, bot *conn.Bot) {
	cId := c.Param("channelId")
	m := util.MustParseReader[apiEntity.Message](c.Request.Body)
	var resp connEntity.SendMessageResponse
	var err error
	var attachmentBytes []byte
	if len(m.Attachments) > 0 {
		attachmentBytes = m.Attachments[0].ToBytes()
	}
	resp, err = bot.PostMessage(m.ToContent(), m.ID, cId, attachmentBytes)
	handleErr(c, err, resp)
}

func createDirectMessage(c *gin.Context, bot *conn.Bot) {
	guildId := c.Param("guildId")
	userId := c.Param("userId")
	message, err := bot.CreateDirectMessage(guildId, userId)
	handleErr(c, err, message)
}

func postDirectMessage(c *gin.Context, bot *conn.Bot) {
	guildId := c.Param("guildId")
	b := util.MustParseReader[dto.MessageToCreate](c.Request.Body)
	directMessage, err := bot.PostDirectMessage(guildId, b)
	handleErr(c, err, directMessage)
}

func deleteMember(c *gin.Context, bot *conn.Bot) {
	guildId := c.Param("guildId")
	userId := c.Param("userId")
	var delDay *int = nil
	var addBlackList *bool = nil
	deleteHistoryMsgDayStr := c.Query("deleteHistoryMsgDay")
	if day, err := strconv.Atoi(deleteHistoryMsgDayStr); err != nil {
		delDay = &day
	}
	addBlackListStr := c.Query("addBlackList")
	if addBlackListStr == "true" {
		var b = true
		addBlackList = &b
	} else if addBlackListStr == "false" {
		var b = false
		addBlackList = &b
	}
	err := bot.DeleteGuildMember(guildId, userId, delDay, addBlackList)
	handleErr(c, err, nil)
}

func getPagedMember(c *gin.Context, bot *conn.Bot) {
	guildId := c.Param("guildId")
	var after string
	var err error
	if after = c.Query("after"); after == "" {
		after = "0"
	}
	var limit int
	if limit, err = strconv.Atoi(c.Query("limit")); err != nil || limit <= 0 || limit > 1000 {
		limit = 400
	}
	mList, err := bot.GuildMembers(guildId, after, limit)
	handleErr(c, err, mList)
}

func delMsg(c *gin.Context, bot *conn.Bot) {
	channelId := c.Param("channelId")
	messageId := c.Param("messageId")
	err := bot.DelMsg(channelId, messageId, c.Query("hideTip") == "true")
	handleErr(c, err, nil)
}
