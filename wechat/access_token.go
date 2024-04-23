package wechat

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const wechatTokenUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET"

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Expire      time.Time
	Code        int    `json:"errcode"`
	Message     string `json:"errmsg"`
}

var accessToken = Token{}

func AccessToken() string {
	if accessToken.Expire.After(time.Now()) {
		return accessToken.AccessToken
	}
	r := strings.ReplaceAll(wechatTokenUrl, "APPID", os.Getenv("yzk_wechat_appId"))
	r = strings.ReplaceAll(r, "APPSECRET", os.Getenv("yzk_wechat_appSecret"))
	response, err := http.Get(r)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	_ = json.Unmarshal(responseData, &accessToken)
	if accessToken.Code == 0 && accessToken.ExpiresIn > 0 {
		accessToken.Expire = time.Now().Add((time.Duration(accessToken.ExpiresIn) - 60) * time.Second)
		return accessToken.AccessToken
	}
	return ""
}
