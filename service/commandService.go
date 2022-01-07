package service

import (
	"fmt"
	"go-bot/utils/animeSearch"
	"go-bot/utils/bilibili"
	"go-bot/utils/fflogs"
	"go-bot/utils/picSearch"
	"go-bot/utils/taro"
	"go-bot/utils/tools"
	"strings"
)

func CommandHandler(data map[string]interface{}) {
	sendData := make(map[string]string)
	if data["message_type"] == "private" {
		sendData["action"] = "send_private_msg"
		sendData["type"] = "user_id"
		sendData["type_id"] = tools.GetString(data["user_id"])
	} else if data["message_type"] == "group" {
		sendData["action"] = "send_group_msg"
		sendData["type"] = "group_id"
		sendData["type_id"] = tools.GetString(data["group_id"])
	}
	sendData["message"] = GetMessage(tools.GetString(data["message"]))
	if sendData["message"] == "" {
		return
	}
	fmt.Println("[发送消息体]:", sendData)
	SendRequest(sendData)
}

// GetMessage 返回信息
func GetMessage(msg string) string {
	str := ""
	//str := "测" + url.QueryEscape("\n") +
	//	"试" +
	//	url.QueryEscape("\n") +
	//	"[CQ:image,file=https://gchat.qpic.cn/gchatpic_new/283213563/3920014266-3050691720-157FB307D61F908C2F5C29F095BD74A6/0?term=3]"

	// 塔罗占卜
	if msg == "taro" {
		str = taro.GetTaro()
	}
	// 查询logs
	if msg == "logs" {
		fflogs.GetRanking(make(map[string]string))
	}
	// 搜图
	if strings.Index(msg, "zpic") >= 0 {
		url := strings.Split(strings.Split(msg, " ")[1], ",")
		picUrl := strings.Replace(strings.Replace(url[2], "url=", "", -1), "]", "", -1)
		str = picSearch.GetPic(picUrl)
	}
	// 搜番
	if strings.Index(msg, "zanime") >= 0 {
		url := strings.Split(strings.Split(msg, " ")[1], ",")
		picUrl := strings.Replace(strings.Replace(url[2], "url=", "", -1), "]", "", -1)
		str = animeSearch.AnimeSearch(picUrl)
	}
	// 主播开播状态
	//if strings.Index(msg, "zlive") >= 0 {
	//	msg := bilibili.GetLiveStatus()
	//	for _, val := range msg {
	//		return val
	//	}
	//}
	// 新增关注的主播
	if strings.Index(msg, "addliver") >= 0 {
		mid := strings.Split(msg, " ")[1]
		str = bilibili.AddLiver(mid)
	}
	return str
}
