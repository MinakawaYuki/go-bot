package service

import (
	"fmt"
	"net/http"
)

// SendRequest 处理私聊请求
func SendRequest(data map[string]string) {
	url := "http://127.0.0.1:5700/" + data["action"] + "?" + data["type"] + "=" + data["type_id"] + "&message=" + data["message"]
	_, err := http.Get(url)
	if err != nil {
		fmt.Println("sendErr", err)
	}
}
