package taro

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
)

var list = []string{
	"愚者",
	"魔术师",
	"女祭司",
	"女皇",
	"皇帝",
	"教皇",
	"恋人",
	"战车",
	"力量",
	"隐者",
	"正义",
	"倒吊者",
	"死神",
	"节制",
	"恶魔",
	"塔",
	"星星",
	"月亮",
	"太阳",
	"审判",
}

var cardType = []string{
	"正位",
	"逆位",
}

func GetTaro() string {
	rand.Seed(time.Now().Unix())
	carNum := rand.Intn(20)
	typeNum := rand.Intn(2)

	card := list[carNum]
	cardState := cardType[typeNum]

	//读取本地json文件
	jsonFile, err := os.Open("utils/taro/taro.json")
	if err != nil {
		fmt.Println("读取文件出错:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json := string(byteValue)

	strCard := card
	strState := gjson.Get(json, card+"."+cardState).String()
	image := gjson.Get(json, card+".image").String()
	cardDetail := gjson.Get(json, card+".牌面解读").String()

	return strings.Replace(strCard+url.QueryEscape("\n")+cardDetail+url.QueryEscape("\n")+cardState+":"+url.QueryEscape("\n")+strState+url.QueryEscape("\n")+"[CQ:image,file="+image+"]", "\n", url.QueryEscape("\n"), -1)
}
