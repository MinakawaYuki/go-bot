package bilibili

import (
	"bytes"
	"context"
	"fmt"
	"go-bot/setting"
	"go-bot/utils/tools"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var livers = []string{"6713974", "190331", "1288041386"}
var liver = []string{"6713974"}

// GetLiveStatus 单次获取某主播开播状态
func GetLiveStatus() []string {
	var msg []string
	for _, val := range liver {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.bilibili.com/x/space/acc/info?mid="+val, nil)
		if err != nil {
			fmt.Println("req Err", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("resp Err", err)
		}
		defer resp.Body.Close()

		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("bodyText Err", err)
		}
		response := tools.Bytes2Map(bodyText)
		fmt.Println("[liver response]:", response)
		response = response["data"].(map[string]interface{})
		up := response["name"].(string)
		liveInfo := response["live_room"].(map[string]interface{})
		if liveInfo["liveStatus"].(float64) == 1 {
			message := up + "开播了\n" +
				"标题:" + liveInfo["title"].(string) + "\n" +
				"直播间地址:" + liveInfo["url"].(string) + "\n" +
				"封面:" + "[CQ:image,file=" + liveInfo["cover"].(string) + "]"
			msg = append(msg, url.QueryEscape(message))
		} else {
			message := up + "是懒狗，根本不播！"
			msg = append(msg, url.QueryEscape(message))
		}
	}
	return msg
}

// GetLiveStatusPerMin 轮询获取主播开播状态
func GetLiveStatusPerMin() {
	for {
		var post = map[string]string{
			"action":  "send_private_msg",
			"type":    "user_id",
			"type_id": "283213563",
			"message": "",
		}
		for _, val := range livers {
			// 先检查redis中是否有数据
			if getExistByRedis(val) == true {
				continue
			}
			client := &http.Client{}
			req, err := http.NewRequest("GET", "https://api.bilibili.com/x/space/acc/info?mid="+val, nil)
			if err != nil {
				fmt.Println("req Err", err)
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("resp Err", err)
			}
			bodyText, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("bodyText Err", err)
			}
			resp.Body.Close()

			response := tools.Bytes2Map(bodyText)
			response = response["data"].(map[string]interface{})
			up := response["name"].(string)
			fmt.Println("[liver status info - " + time.Now().Format("2006/1/02 15:04") + "]:正在查询" + up + "的直播状态....")
			liveInfo := response["live_room"].(map[string]interface{})
			if liveInfo["liveStatus"].(float64) == 1 {
				fmt.Println("[liver status info - " + time.Now().Format("2006/1/02 15:04") + "]:查询结果为正在直播....")
				//修改redis中的主播直播间状态
				if getStatusByRedis(val) == true {
					err := setStatusByRedis(val, liveInfo["liveStatus"].(float64))
					if err != nil {
						fmt.Println("[redis]保存直播状态"+liveInfo["liveStatus"].(string)+"失败:", err)
					}
				} else {
					message := up + "开播了\n" +
						"标题:" + liveInfo["title"].(string) + "\n" +
						"直播间地址:" + liveInfo["url"].(string) + "\n" +
						"封面:" + "[CQ:image,file=" + liveInfo["cover"].(string) + "]"
					post["message"] = url.QueryEscape(message)
					sendRequest(post)
				}
			} else {
				fmt.Println("[liver status info - " + time.Now().Format("2006/1/02 15:04") + "]:查询结果为未开播....")
				//修改redis中的主播直播间状态
				if getStatusByRedis(val) == false {
					err := setStatusByRedis(val, liveInfo["liveStatus"].(float64))
					if err != nil {
						fmt.Println("[redis]保存直播状态失败:", err)
					}
				} else {
					message := up + "下播了\n"
					post["message"] = url.QueryEscape(message)
					sendRequest(post)
				}
			}
		}
		time.Sleep(time.Second * 150)
	}
}

// getStatusByRedis 查询redis中的直播间状态
func getStatusByRedis(mid string) bool {
	ctx := context.Background()
	client := setting.RedisClient
	status, _ := client.Get(ctx, mid).Float64()
	if status == 1 {
		return true
	}
	return false
}

// setStatusByRedis 保存直播间状态
func setStatusByRedis(mid string, status float64) error {
	ctx := context.Background()
	client := setting.RedisClient
	_, err := client.Set(ctx, mid, status, time.Minute).Result()
	return err
}

// getExistByRedis 查询key是否存在
func getExistByRedis(mid string) bool {
	ctx := context.Background()
	status := setting.RedisClient.Exists(ctx, mid).Val()
	if status == 1 {
		return true
	}
	return false
}

// sendRequest 发送消息
func sendRequest(data map[string]string) {
	sendUrl := "http://" + setting.BotSetting.IP + ":" + setting.BotSetting.Port + "/" + data["action"] + "?" + data["type"] + "=" + data["type_id"] + "&message=" + data["message"]
	client := &http.Client{}
	req, err := http.NewRequest("GET", sendUrl, nil)
	if err != nil {
		fmt.Println("sendErr", err)
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	fmt.Println("[resp]", bytes.NewBuffer(bodyText).String())
}
