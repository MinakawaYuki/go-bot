package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-bot/service"
	"go-bot/utils/tools"
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
		data := tools.Bytes2Map(message)
		if data["meta_event_type"] != "heartbeat" {
			if data["post_type"] == "message" && (data["message_type"] == "private" || data["message_type"] == "group") {
				fmt.Println("[收到message]:", tools.Bytes2Map(message))
				service.CommandHandler(data)
			}
		}
		if err != nil {
			break
		}
	}
}
