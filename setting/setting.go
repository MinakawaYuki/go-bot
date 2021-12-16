package setting

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gopkg.in/ini.v1"
	"time"
)

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	DB          int
}

type Bot struct {
	IP   string
	Port string
}

var RedisSetting = &Redis{}
var RedisClient *redis.Client

var BotSetting = &Bot{}

func SetUp() {
	Cfg, err := ini.Load("conf/config.ini")
	if err != nil {
		fmt.Println("Fail to parse 'conf/app.ini': ", err)
	}
	err = Cfg.Section("redis").MapTo(&RedisSetting)
	if err != nil {
		fmt.Println("Cfg.MapTo RedisSetting err: ", err)
	}
	err = Cfg.Section("bot").MapTo(&BotSetting)
	if err != nil {
		fmt.Println("Cfg.MapTo BotSetting err: ", err)
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     RedisSetting.Host,
		Password: RedisSetting.Password,
		DB:       RedisSetting.DB})
}
