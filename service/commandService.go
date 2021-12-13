package service

import (
	"fmt"
	"go-bot/utils/fflogs"
	"go-bot/utils/taro"
	"math/big"
)

func CommandHandler(data map[string]interface{}) {
	sendData := make(map[string]string)
	if data["message_type"] == "private" {
		sendData["action"] = "send_private_msg"
		sendData["type"] = "user_id"
		sendData["type_id"] = GetString(data["user_id"])
	} else if data["message_type"] == "group" {
		sendData["action"] = "send_group_msg"
		sendData["type"] = "group_id"
		sendData["type_id"] = GetString(data["group_id"])
	}
	sendData["message"] = GetMessage(GetString(data["raw_message"]))
	if sendData["message"] == "" {
		return
	}
	fmt.Println("[发送消息体]:", sendData)
	SendRequest(sendData)
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
	case string:
		str = v.(string)
		break
	}
	return str
}

// GetMessage 返回信息
func GetMessage(msg string) string {
	str := ""
	//str := "测" + url.QueryEscape("\n") +
	//	"试" +
	//	url.QueryEscape("\n") +
	//	"[CQ:image,file=https://gchat.qpic.cn/gchatpic_new/283213563/3920014266-3050691720-157FB307D61F908C2F5C29F095BD74A6/0?term=3]"
	if msg == "taro" {
		str = taro.GetTaro()
	}
	if msg == "logs" {
		fflogs.GetRanking(make(map[string]string))
	}
	return str
}
