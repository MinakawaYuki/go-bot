package fflogs

import "fmt"

// V2 v2接口参数 目前不会graphQL 暂时搁置
var V2 = map[string]string{
	"ClientId":     "9513e857-07ee-4130-88f7-18857a58d291",
	"ClientSecret": "arIyugOHwyHJcgsymLjjvpz6k4R7iUWNo1QdzfFH",
	"OpenApi":      "https://www.fflogs.com/api/v2/client",
}

// V1 v1版本接口key
var V1 = map[string]string{
	"Key": "ad47257a0254969cd179b849624920f7",
}

type Client struct {
	ClientId     string
	ClientSecret string
}

func GetRanking(data map[string]string) map[string]string {
	accessToken := GetAccessToken(Client{
		ClientId:     V2["ClientId"],
		ClientSecret: V2["ClientSecret"],
	})
	fmt.Println("accessToken:", accessToken)
	return make(map[string]string)
}
