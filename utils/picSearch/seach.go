package picSearch

import (
	"go-bot/utils/tools"
	"io/ioutil"
	"log"
	"net/http"
	urrl "net/url"
	"strconv"
)

var apiKey = "0a68561f665006f23ef087475b5d9ce861548129"
var url = "https://saucenao.com/search.php?testmode=1&output_type=2&api_key=" + apiKey + "&numres=5&db=999&url="
var split = urrl.QueryEscape("\n")

func GetPic(picUrl string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url+picUrl, nil)
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
	results := response["results"]
	message := ""
	for key, val := range results.([]interface{}) {
		header := val.(map[string]interface{})["header"].(map[string]interface{})
		data := val.(map[string]interface{})["data"].(map[string]interface{})

		if key > 1 {
			break
		}
		message += "第" + strconv.Itoa(key+1) + "个结果:" + split +
			//"缩略图:" + header["thumbnail"].(string) + split +
			getImageInfo(data, header["index_id"])
	}
	return message
}

// getImageInfo 拼接返回图片的具体信息
func getImageInfo(data map[string]interface{}, indexId interface{}) string {
	msg := ""
	switch indexId.(float64) {
	case 40: //Index #40: FurAffinity
		msg += "标题:" + data["title"].(string) + split + "作者:" + data["author_name"].(string) + split + "作者主页:" + data["author_url"].(string) + split + "链接:" + data["ext_urls"].([]interface{})[0].(string)
		break
	case 5: //Index #5: Pixiv Images
		msg += "标题:" + data["title"].(string) + split + "作者:" + data["member_name"].(string) + split + "链接:" + data["ext_urls"].([]interface{})[0].(string)
		break
	case 22: //Index #22: H-Anime*
		msg += "source:" + data["source"].(string) + split + "part:" + data["part"].(string) + split + "year:" + data["year"].(string) + split + "链接:" + data["ext_urls"].([]interface{})[0].(string)
		break
	case 18: //Index #18: H-Misc (nhentai)
		msg += "source:" + data["source"].(string) + split + "creator:" + data["creator"].(string) + split + "eng_name:" + data["eng_name"].(string) + split + "jp_name:" + data["jp_name"].(string)
		break
	case 37: //Index #37: MangaDex
		msg += "source:" + data["source"].(string) + split + "artist:" + data["artist"].(string) + split + "链接:" + data["ext_urls"].([]interface{})[0].(string)
		break
	case 25: //Index #25: Gelbooru
		msg += "material:" + data["material"].(string) + split + "creator:" + data["creator"].(string) + split + "链接:" + data["ext_urls"].([]interface{})[0].(string)
		break
	case 38: //Index #38: H-Misc (E-Hentai)
		msg += "source:" + data["source"].(string) + split + "creator:" + data["creator"].(string) + split + "eng_name:" + data["eng_name"].(string) + split + "jp_name:" + data["jp_name"].(string)
		break
	case 36: //Index #36: Madokami (Manga)
		msg += "source:" + data["source"].(string) + split + "part:" + data["part"].(string) + split + "type:" + data["type"].(string)
		break
	case 21: //Index #21: Anime
		msg += "source:" + data["source"].(string) + split + "part:" + data["part"].(string) + split + "year:" + data["year"].(string)
		break
	case 41: //Index #41: Twitter
		msg += "推文:" + data["tweet_id"].(string) + split + "推主:" + data["twitter_user_id"].(string) + split + "链接:" + data["ext_urls"].([]interface{})[0].(string)
		break
	case 8: //Index #8: Nico Nico Seiga
		msg += "title:" + data["title"].(string) + split + "作者:" + data["member_name"].(string) + split + "链接:" + data["ext_urls"].([]interface{})[0].(string)
		break
	case 32: //Index #32: bcy.net Cosplay
		msg += "title:" + data["title"].(string) + split + "作者:" + data["member_name"].(string) + split + "链接:" + data["ext_urls"].([]interface{})[0].(string)
		break
	default:

	}
	return msg
}
