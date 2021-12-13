package main

import (
	"go-bot/router"
	"go-bot/setting"
)

func main() {
	setting.SetUp()
	router.InitRouter()
}
