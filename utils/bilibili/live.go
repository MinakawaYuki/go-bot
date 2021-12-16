package bilibili

import (
	"fmt"
	"go-bot/utils/tools"
	"io/ioutil"
	"net/http"
	"net/url"
)

var livers = []string{"6713974"}

func GetLiveStatus() []string {
	var msg []string
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
		}
	}
	return msg
}
