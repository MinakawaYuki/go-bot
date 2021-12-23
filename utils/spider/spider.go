package spider

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"go-bot/setting"
)

var Spider = colly.NewCollector(
	colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"),
	colly.MaxDepth(1),
	//colly.Debugger(&debug.LogDebugger{}),
)
var SUrl = "https://www.jianshu.com"

func Result() {
	Spider.OnHTML("ul[class='note-list']", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, element *colly.HTMLElement) {
			//href
			href := element.ChildAttrs("div[class='content'] > a[class='title']", "href")
			//title
			title := element.ChildText("div[class='content'] > a[class='title']")
			//summary
			summary := element.ChildText("div[class='content'] > p[class='abstract']")
			fmt.Println("title:", title)
			fmt.Println("href:", href)
			fmt.Println("summary:", summary)
		})
	})
	err := Spider.Visit(SUrl)
	if err != nil {
		setting.Log.Fatal("[spider]:", err)
	}
}
