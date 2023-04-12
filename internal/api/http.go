package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tencent-connect/botgo/dto"
	"html/template"
	"net/http"
	apiEntity "qq-guild-bot/internal/api/entity"
	"qq-guild-bot/internal/conn"
	connEntity "qq-guild-bot/internal/conn/entity"
	"qq-guild-bot/internal/embeded"
	myLogger "qq-guild-bot/internal/pkg/log"
	"qq-guild-bot/internal/pkg/util"
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
	guildApi.GET("/member", getAllMember)
	guildApi.PATCH("/member/:userId", updateUser)
	guildApi.DELETE("/member", banByBatch)
	guildApi.PUT("/direct_message/:userId", createDirectMessage)
	guildApi.POST("/direct_message", postDirectMsg)

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

func updateUser(c *gin.Context) {
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
	resp, err = conn.PostMsg(m.ToContent(), m.ID, cId, attachmentBytes)
	handleErr(c, err, resp)
}

func createDirectMessage(c *gin.Context) {
	guildId := c.Param("guildId")
	userId := c.Param("userId")
	message, err := conn.CreateDirectMessage(guildId, userId)
	handleErr(c, err, message)
}

func postDirectMsg(c *gin.Context) {
	guildId := c.Param("guildId")
	b := util.MustParseReader[dto.MessageToCreate](c.Request.Body)
	directMessage, err := conn.PostDirectMessage(guildId, b)
	handleErr(c, err, directMessage)
}

func banByBatch(c *gin.Context) {
	guildId := c.Param("guildId")
	memberIdList := util.MustParseArrayReader[string](c.Request.Body)
	err := conn.BanByBatch(guildId, memberIdList)
	handleErr(c, err, nil)
}

func getAllMember(c *gin.Context) {
	guildId := c.Param("guildId")
	mList := conn.GetAllMember(guildId)
	c.JSON(http.StatusOK, apiEntity.NewOkResp(mList, nil))
}

func delMsg(c *gin.Context) {
	channelId := c.Param("channelId")
	messageId := c.Param("messageId")
	err := conn.DelMsg(channelId, messageId, c.Query("hideTip") == "true")
	handleErr(c, err, nil)
}
