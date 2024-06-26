package main

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"sort"
	"strings"
	"time"
)

func PrintRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// t := time.Now()
		// 请求前
		dumpRequest, _ := httputil.DumpRequest(c.Request, true)
		log.Println("Request:\n" + string(dumpRequest))
		c.Next()
		// 请求后
		// latency := time.Since(t)
		// log.Print("cost time: " + latency.String())
		// status := c.Writer.Status()
		// log.Println(status)
	}
}
func checkSignature(signature string, timestamp string, nonce string) bool {
	token := "yzk"
	tmpArr := []string{token, timestamp, nonce}
	sort.Strings(tmpArr)
	tmpStr := strings.Join(tmpArr, "")
	hash := sha1.New()
	hash.Write([]byte(tmpStr))
	tmpHash := fmt.Sprintf("%x", hash.Sum(nil))
	return tmpHash == signature
}

type TextMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	Image        Image    `xml:"Image"`
	MsgId        string   `xml:"MsgId"`
}
type Image struct {
	MediaId string `xml:"MediaId"`
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})
	// r.Use(PrintRequest())
	r.GET("/yzk/wechat/notify/:appId", func(c *gin.Context) {
		log.Printf("appId: %s\n", c.Param("appId"))
		echostr := c.Query("echostr")
		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		if checkSignature(signature, timestamp, nonce) {
			c.String(http.StatusOK, echostr)
		}
	})
	r.POST("/yzk/wechat/notify/:appId", func(c *gin.Context) {
		msg := TextMsg{}
		_ = c.ShouldBindXML(&msg)
		log.Printf("appId[%s]openId[%s]MsgType[%s]Content[%s]MsgId[%s]\n",
			c.Param("appId"), msg.FromUserName, msg.MsgType, msg.Content, msg.MsgId)
		if "text" != msg.MsgType {
			c.String(http.StatusOK, "success")
			return
		}
		f := msg.FromUserName
		msg.FromUserName = msg.ToUserName
		msg.ToUserName = f
		msg.CreateTime = time.Now().Unix()
		msg.MsgType = "image"
		msg.Image.MediaId = "4B_rIg-H5pEknQ0hyPf-TsRBBrbSvWY_wHOM-WycMrUa9Ix8DJzahYUjaaRJCdKe"
		c.XML(http.StatusOK, msg)
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		c.JSON(http.StatusOK, gin.H{"user": user})
	})
	return r
}

func main() {
	_ = setupRouter().Run("127.0.0.1:8351")
}
