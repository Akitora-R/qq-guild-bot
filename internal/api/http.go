package api

import (
	"encoding/json"
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
		c.JSON(200, conn.GetSelf())
	})
	guildApi := group.Group("/guild/:guildId")
	guildApi.GET("/member", getPagedMember)
	guildApi.GET("/member/:userId", getMemberDetail)
	guildApi.PATCH("/member/:userId", updateMember)
	guildApi.PUT("/member/:userId/roles/:roleId", updateMemberRole)
	guildApi.DELETE("/member/:userId/roles/:roleId", deleteMemberRole)
	guildApi.DELETE("/member/:userId", deleteMember)
	guildApi.PUT("/direct-message/:userId", createDirectMessage)
	guildApi.POST("/direct-message", postDirectMessage)
	guildApi.GET("/roles", getRoles)

	channelApi := group.Group("/channel/:channelId")
	channelApi.POST("/message", sendMsg)
	channelApi.DELETE("/message/:messageId", delMsg)
	if err := engine.Run("127.0.0.1:6800"); err != nil {
		panic(err)
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

func getRoles(c *gin.Context) {
	roles, err := conn.GetRoles(c.Param("guildId"))
	handleErr(c, err, roles)
}

func getMemberDetail(c *gin.Context) {
	guildId := c.Param("guildId")
	userId := c.Param("userId")
	member, err := conn.GetGuildMember(guildId, userId)
	handleErr(c, err, member)
}

func updateMemberRole(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	var b dto.MemberAddRoleBody
	if err == nil {
		_ = json.Unmarshal(bodyBytes, &b)
	}
	err = conn.AddMemberRole(c.Param("guildId"), c.Param("roleId"), c.Param("userId"), &b)
	handleErr(c, err, nil)
}

func deleteMemberRole(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	var b dto.MemberAddRoleBody
	if err == nil {
		_ = json.Unmarshal(bodyBytes, &b)
	}
	err = conn.DeleteMemberRole(c.Param("guildId"), c.Param("roleId"), c.Param("userId"), &b)
	handleErr(c, err, nil)
}

func updateMember(c *gin.Context) {
	u := util.MustParseReader[connEntity.UserUpdate](c.Request.Body)
	guildId := c.Param("guildId")
	userId := c.Param("userId")
	handleErr(c, conn.UpdateUser(u, guildId, userId), nil)
}

func sendMsg(c *gin.Context) {
	cId := c.Param("channelId")
	m := util.MustParseReader[apiEntity.Message](c.Request.Body)
	var resp connEntity.SendMessageResponse
	var err error
	var attachmentBytes []byte
	if len(m.Attachments) > 0 {
		attachmentBytes = m.Attachments[0].ToBytes()
	}
	resp, err = conn.PostMessage(m.ToContent(), m.ID, cId, attachmentBytes)
	handleErr(c, err, resp)
}

func createDirectMessage(c *gin.Context) {
	guildId := c.Param("guildId")
	userId := c.Param("userId")
	message, err := conn.CreateDirectMessage(guildId, userId)
	handleErr(c, err, message)
}

func postDirectMessage(c *gin.Context) {
	guildId := c.Param("guildId")
	b := util.MustParseReader[dto.MessageToCreate](c.Request.Body)
	directMessage, err := conn.PostDirectMessage(guildId, b)
	handleErr(c, err, directMessage)
}

func deleteMember(c *gin.Context) {
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
	err := conn.DeleteGuildMember(guildId, userId, delDay, addBlackList)
	handleErr(c, err, nil)
}

func getPagedMember(c *gin.Context) {
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
	mList, err := conn.GuildMembers(guildId, after, limit)
	handleErr(c, err, mList)
}

func delMsg(c *gin.Context) {
	channelId := c.Param("channelId")
	messageId := c.Param("messageId")
	err := conn.DelMsg(channelId, messageId, c.Query("hideTip") == "true")
	handleErr(c, err, nil)
}
