package main

import (
	"go-bot/router"
	"go-bot/setting"
	"go-bot/utils/bilibili"
)

func main() {
	setting.SetUp()
	//bilibili.WordCloud()
	if setting.PluginSetting.Live == true {
		go bilibili.GetLiveStatusPerMin()
	}
	if setting.PluginSetting.Danmaku == true {
		go bilibili.GetDanmaku()
	}
	router.InitRouter()
}
