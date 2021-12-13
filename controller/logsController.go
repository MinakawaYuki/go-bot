package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"go-bot/service"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// GetEventMessage 接受信息
func GetEventMessage(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		data := Bytes2Map(message)
		if data["meta_event_type"] != "heartbeat" {
			if data["post_type"] == "message" && (data["message_type"] == "private" || data["message_type"] == "group") {
				fmt.Println("[收到message]:", Bytes2Map(message))
				service.CommandHandler(data)
			}
		}
		if err != nil {
			break
		}
	}
}

// Bytes2Map 将ws获取的字节流转换为map并返回
func Bytes2Map(p []byte) (data map[string]interface{}) {
	json := bytes.NewBuffer(p).String()
	data, err := gjson.Parse(json).Value().(map[string]interface{})
	if !err {
		panic("解析json错误")
	}

	return data
}
