package main

import (
	"go-bot/router"
	"go-bot/setting"
	"go-bot/utils/bilibili"
)

func main() {
	setting.SetUp()
	go bilibili.GetLiveStatusPerMin()
	router.InitRouter()
}
