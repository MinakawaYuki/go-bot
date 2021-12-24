package bilibili

import (
	"errors"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/jinzhu/gorm"
	"go-bot/setting"
	"go-bot/utils/tools"
	"io/ioutil"
	"net/http"
	"os"
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

type danmus struct {
	Text  string
	Count int
}

func GetDanmaku() {
	//for GetStatusByRedis(mid) == true {
	for {
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

		if response == nil {
			continue
		}
		room := response.(map[string]interface{})["room"]

		danmakuList := room.([]interface{})
		//length := len(danmakuList)
		//data := danmakuList[length-1].(map[string]interface{})

		for _, item := range danmakuList {
			data := item.(map[string]interface{})
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
		}

		//休眠
		time.Sleep(time.Millisecond * 500)
	}
	//fmt.Println("[主播下播了---close 数据库连接]")
	//defer setting.Db.Close()
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

func WordCloud() {
	var list []danmus
	setting.Db.Raw("SELECT `text`,count( `text` ) AS count FROM `danmaku` where `text` not in ('?','？？？','？？','？？？？','？？？？？','？？？？？？','？？？？？？？','???') GROUP BY `text` HAVING COUNT(`count`) > 10 ORDER BY `count` DESC").Scan(&list)
	var items = make([]opts.WordCloudData, 0)
	for _, v := range list {
		items = append(items, opts.WordCloudData{Name: v.Text, Value: v.Count})
	}
	wc := charts.NewWordCloud()
	wc.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: time.Now().Format("2006-1-02") + " 词云"}))
	wc.AddSeries("wordcloud", items).SetSeriesOptions(charts.WithWorldCloudChartOpts(opts.WordCloudChart{SizeRange: []float32{40, 80}, Shape: "cardioid"}))

	f, _ := os.Create("wordCloud.html")
	_ = wc.Render(f)
}
