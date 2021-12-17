package main

import (
	"go-bot/router"
	"go-bot/setting"
	"go-bot/utils/bilibili"
)

func main() {
	setting.SetUp()
	//bilibili.WordCloud()
	go bilibili.GetLiveStatusPerMin()
	go bilibili.GetDanmaku()
	router.InitRouter()
}
