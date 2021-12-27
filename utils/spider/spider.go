package spider

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"go-bot/setting"
	"time"
)

/**
gbf wiki
标签<div class="ability-table-row ec-1-m" id=".E6.8A.80.E8.83.BD1">技能1
标签<div class="ability-table-row ec-1-m" id=".E6.8A.80.E8.83.BD2">技能2
标签<div class="ability-table-row ec-1-m" id=".E6.8A.80.E8.83.BD3">技能3
标签<div class="ability-table-row ec-1-m" id=".E6.8A.80.E8.83.BD4">技能4
标签<div class="ability-table-row ec-1-m" id=".E5.A5.A5.E4.B9.891">未终突奥义
标签<div class="ability-table-row ec-1-m" id=".E5.A5.A5.E4.B9.892">终突奥义
标签<div class="ability-table-row ec-1-m" id=".E5.A5.A5.E4.B9.893">超限奥义 <div class="ability-desc ab-block full-size">…</div> 技能效果
*/

// Character 角色
type Character struct {
	Name         string `json:"name,omitempty"`
	HP           string `json:"HP"`
	Attack       string `json:"attack" gorm:"comment:'攻击力'"`
	Ougi         string `json:"ougi" gorm:"column:ougi;type:longtext;comment:'奥义'"`
	Ougi2        string `json:"ougi2" gorm:"column:ougi2;type:longtext;comment:'终突奥义'"`
	Ougi3        string `json:"ougi3" gorm:"column:ougi3;type:longtext;comment:'超限奥义'"`
	Ability      string `json:"ability" gorm:"column:ability;type:longtext;comment:'技能'"`
	PassiveSkill string `json:"passive_skill" gorm:"column:passive_skill;type:longtext;comment:'被动技能'"`
	LBSkill      string `json:"LBSkill" gorm:"column:LBSkill;type:longtext;comment:'LB技能'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Skill 技能
type Skill struct {
	Name      string
	CountDown float64
	Detail    string
}

// PassiveSkill 被动技能
type PassiveSkill struct {
	Name   string
	Detail string
}

// LBSkill 被动技能
type LBSkill struct {
	Name   string
	Detail string
}

var Spider = colly.NewCollector(
	colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"),
	colly.MaxDepth(1),
	//colly.Debugger(&debug.LogDebugger{}),
)
var SUrl = "https://gbf.huijiwiki.com/wiki/Char/3040035000"

func Scrape() {
	Spider.OnHTML("div[class='ability-table-row ec-1-m']", func(e *colly.HTMLElement) {
		//test := e.DOM.Find("div[class='name title-text ec-2-bg-m']  > span[class='label-item label-color-purple learn-lv']")
		//fmt.Println(test.Text())
		//e.ForEach("div[class='ability-name']", func(i int, element *colly.HTMLElement) {
		//	text := element.ChildText("div[class='name title-text ec-2-bg-m'] > span[class='name_jp jp']")
		//	fmt.Println(text)
		//})
		e.ForEach("div[class='ability-name']", func(i int, element *colly.HTMLElement) {
			fmt.Println(element.Text)
		})
	})
	err := Spider.Visit(SUrl)
	if err != nil {
		setting.Log.Fatal("[spider]:", err)
	}
}
