package controller

import (
	"fmt"
	"net/http"
)

// HandlePrivateRequest 处理私聊请求
func HandlePrivateRequest(data map[string]string) {
	url := "http://127.0.0.1:5700/" + data["action"] + "?user_id=" + data["user_id"] + "&message=" + data["message"]
	_, err := http.Get(url)
	if err != nil {
		fmt.Println("sendErr", err)
	}
}

// HandleGroupRequest 处理私聊请求
func HandleGroupRequest(data map[string]string) {
	url := "http://127.0.0.1:5700/" + data["action"] + "?user_id=" + data["user_id"] + "&message=" + data["message"]
	_, err := http.Get(url)
	if err != nil {
		fmt.Println("sendErr", err)
	}
}
