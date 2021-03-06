package bilibili

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"go-bot/setting"
	"go-bot/utils/tools"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var livers = []string{"6713974", "190331", "1934433723"}

// Liver 主播
type Liver struct {
	Nickname   string  `json:"nickname,omitempty"`
	UserId     string  `json:"user_id,omitempty"`
	RoomId     float64 `json:"room_id,omitempty"`
	LiveStatus float64 `json:"live_status,omitempty"`
	RoomUrl    string  `json:"room_url,omitempty"`
	Title      string  `json:"title,omitempty"`
	Cover      string  `json:"cover,omitempty"`
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
			if response != nil {
				response = response["data"].(map[string]interface{})
				up := response["name"].(string)
				fmt.Println("[liver status info - " + time.Now().Format("2006/1/02 15:04") + "]:正在查询" + up + "的直播状态....")
				liveInfo, err := response["live_room"].(map[string]interface{})
				if err != true {
					message := "断言失败"
					post["message"] = url.QueryEscape(message)
					sendRequest(post)
					continue
				}
				//使用数据库判定是否开播 启用redis
				var liver Liver
				setting.Db.Where("user_id = ? ", val).First(&liver)
				if liveInfo["liveStatus"].(float64) != liver.LiveStatus {
					if liver.LiveStatus == 0 {
						message := up + "开播了\n" +
							"标题:" + liveInfo["title"].(string) + "\n" +
							"直播间地址:" + liveInfo["url"].(string) + "\n" +
							"封面:" + "[CQ:image,file=" + liveInfo["cover"].(string) + "]"
						post["message"] = url.QueryEscape(message)
						sendRequest(post)
					} else {
						message := up + "下播了\n"
						post["message"] = url.QueryEscape(message)
						sendRequest(post)
					}
					// 保存直播状态
					setting.Db.Table("liver").Where("user_id = ? ", val).Update(map[string]interface{}{"title": liveInfo["title"].(string), "cover": liveInfo["cover"].(string), "live_status": liveInfo["liveStatus"].(float64)})
				}
			}
		}
		time.Sleep(time.Second * 30)
	}
}

// AddLiver 增加关注的主播
func AddLiver(mid string) string {
	var msg string

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/space/acc/info?mid="+mid, nil)
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
	response = response["data"].(map[string]interface{})
	up := response["name"].(string)
	liveInfo := response["live_room"].(map[string]interface{})

	liver := Liver{
		Nickname:   up,
		UserId:     mid,
		RoomId:     liveInfo["roomid"].(float64),
		LiveStatus: liveInfo["liveStatus"].(float64),
		RoomUrl:    liveInfo["url"].(string),
		Title:      liveInfo["title"].(string),
		Cover:      liveInfo["cover"].(string),
	}
	if errors.Is(setting.Db.Where(&Liver{UserId: liver.UserId}).First(&Liver{}).Error, gorm.ErrRecordNotFound) {
		err = setting.Db.Create(&liver).Error
		if err != nil {
			msg = "保存数据出错:" + err.Error()
		}
		msg = "添加成功"
	} else {
		msg = "该数据已存在"
	}
	return msg
}

// GetStatusByRedis 查询redis中的直播间状态
func GetStatusByRedis(mid string) bool {
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
	if status == 1 {
		_, err := client.Set(ctx, mid, status, 0).Result()
		if err != nil {
			return err
		}
	} else {
		_, err := client.Set(ctx, mid, status, 0).Result()
		if err != nil {
			return err
		}
	}
	return nil
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
