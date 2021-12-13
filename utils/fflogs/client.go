package fflogs

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"
	"go-bot/setting"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var ctx = context.Background()

// GetAccessToken 获取token
func GetAccessToken(client Client) string {
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", "https://www.fflogs.com/oauth/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(client.ClientId, client.ClientSecret)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	accessToken := ""
	//todo 超时处理
	// 使用redis存储token 过期时间 半小时
	val, err := setting.RedisClient.Get(ctx, "logs_token").Result()
	if err == redis.Nil {
		//不存在键
		accessToken = gjson.GetBytes(bodyText, "access_token").String()
		setting.RedisClient.Set(ctx, "logs_token", accessToken, 15*time.Minute)
	} else if err != nil {
		fmt.Println("[redis err]:", err)
	} else {
		accessToken = val
	}
	return accessToken
}
