package fflogs

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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
	accessToken := gjson.GetBytes(bodyText, "access_token").String()
	//todo 超时处理
	// 使用redis存储token 比expires_in 短
	//
	return accessToken
}
