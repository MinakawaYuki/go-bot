package animeSearch

import (
	"go-bot/utils/tools"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var searchUrl = "https://api.trace.moe/search?url="
var split = url.QueryEscape("\n")

func AnimeSearch(picUrl string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchUrl+picUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	response := tools.Bytes2Map(bodyText)
	result := response["result"].([]interface{})
	return getAnimeInfo(result)
}

func getAnimeInfo(res []interface{}) string {
	str := ""
	for key, val := range res {
		if key > 0 {
			break
		}
		data := val.(map[string]interface{})
		filename := data["filename"].(string)
		if strings.Index(filename, "[") >= 0 {
			filename = strings.Replace(filename, "[", "", -1)
			filename = strings.Replace(filename, "]", "", -1)
		}

		str += "标题:" + filename + split +
			"相似度:" + tools.GetString(data["similarity"].(float64)) + split +
			"from:" + tools.GetString(data["from"].(float64)) + split +
			"to:" + tools.GetString(data["to"].(float64)) + split +
			"image:" + data["image"].(string)
	}

	return str
}
