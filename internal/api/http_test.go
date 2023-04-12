package api_test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tencent-connect/botgo/log"
	"qq-guild-bot/internal/pkg/config"
	"qq-guild-bot/internal/pkg/util/http"
	"strconv"
	"testing"
)

func TestSendMsg(t *testing.T) {
	authToken := fmt.Sprintf("%s.%s", strconv.FormatUint(config.AppConf.AppID, 10), config.AppConf.AccessToken)
	req := resty.New().R().SetAuthScheme("Bot").SetAuthToken(authToken)
	r, err := req.SetMultipartFormData(map[string]string{
		"content": "content",
		"msg_id":  "08d590bb888addcced850110a9938c01389d0148bbcfbf9b06",
	}).SetPathParam("channel_id", "2296233").Post(config.BaseApi + "/channels/{channel_id}/messages")
	if err != nil {
		panic(err)
	}
	log.Info(string(r.Body()))
}

func TestSendMsg3(t *testing.T) {

}

// sample https://twitter.com/phy_sen/status/1591413038548422656/photo/1
func TestSendMsg2(t *testing.T) {
	authToken := fmt.Sprintf("%s.%s", strconv.FormatUint(config.AppConf.AppID, 10), config.AppConf.AccessToken)
	d := map[string]string{
		"content": "111333",
		"msg_id":  "08d590bb888addcced850110aa938c0138970b48cd99c49b06",
	}
	form, err := http.PostMultipartForm(map[string]string{"Authorization": "Bot " + authToken}, d, nil, config.BaseApi+"/channels/2296234/messages")
	if err != nil {
		panic(err)
	}
	log.Info(string(form))
}
