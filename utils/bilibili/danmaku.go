package bilibili

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"go-bot/setting"
	"go-bot/utils/tools"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var liveUrl = "https://api.live.bilibili.com/xlive/web-room/v1/dM/gethistory"
var roomId = "7129270"
var mid = "6713974"

type danmaku struct {
	Text     string `gorm:"text" json:"text"`
	Uid      string `gorm:"uid" json:"uid"`
	Nickname string `gorm:"nickname" json:"nickname"`
	Timeline string `gorm:"timeline" json:"timeline"`
}

func GetDanmaku() {
	for GetStatusByRedis(mid) == true {
		client := &http.Client{}
		req, err := http.NewRequest("POST", liveUrl, strings.NewReader("roomid="+roomId))
		if err != nil {
			fmt.Println("[get danmaku err]:", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Host", "api.live.bilibili.com")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")

		resq, err := client.Do(req)

		if err != nil {
			fmt.Println("[get danmaku err]:", err)
		}

		body, err := ioutil.ReadAll(resq.Body)
		if err != nil {
			fmt.Println("[get danmaku err]:", err)
		}
		resq.Body.Close()

		response := tools.Bytes2Map(body)["data"]

		room := response.(map[string]interface{})["room"]

		danmakuList := room.([]interface{})
		length := len(danmakuList)
		data := danmakuList[length-1].(map[string]interface{})

		dan := danmaku{
			Text:     data["text"].(string),
			Uid:      fmt.Sprintf("%f", data["uid"].(float64)),
			Nickname: data["nickname"].(string),
			Timeline: data["timeline"].(string),
		}
		if exsit(dan) {
			err := add(dan)
			if err != nil {
				fmt.Println("[add danmaku err]:", err)
			}
		}
		//休眠
		time.Sleep(time.Millisecond * 500)
	}
	defer setting.Db.Close()
}

func exsit(data danmaku) bool {
	if errors.Is(setting.Db.Where(&data).First(&danmaku{}).Error, gorm.ErrRecordNotFound) {
		return true
	}
	return false
}

func add(data danmaku) error {
	return setting.Db.Create(&data).Error
}
