package service

import (
	"bytes"
	"fmt"
	"go-bot/setting"
	"io/ioutil"
	"log"
	"net/http"
)

// SendRequest 处理私聊请求
func SendRequest(data map[string]string) {
	url := "http://" + setting.BotSetting.IP + ":" + setting.BotSetting.Port + "/" + data["action"] + "?" + data["type"] + "=" + data["type_id"] + "&message=" + data["message"]
	//fmt.Println("[url]:", url)
	//resp, err := http.Get(url)
	//if err != nil {
	//	fmt.Println("sendErr", err)
	//}
	//fmt.Println("[resp]", resp)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
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
