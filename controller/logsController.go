package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"math/big"
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
			fmt.Println("buf", Bytes2Map(message))
			if data["message_type"] == "private" && data["post_type"] == "message" && data["raw_message"] == "测试" {
				sendData := map[string]string{
					"action":  "send_private_msg",
					"user_id": GetString(data["user_id"]),
					"message": "回复",
				}
				HandlePrivateRequest(sendData)
			}
		}
		if err != nil {
			break
		}
	}
}

// GetString 类型断言
func GetString(v interface{}) string {
	var str string
	switch v.(type) {
	case float64:
		vv := v.(float64)
		data := big.NewRat(1, 1)
		data.SetFloat64(vv)
		str = data.FloatString(0)
		break
	}
	println("str", str)
	return str
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
