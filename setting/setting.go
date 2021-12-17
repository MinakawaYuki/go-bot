package setting

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/ini.v1"
	"log"
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

type Mysql struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

var RedisSetting = &Redis{}
var RedisClient *redis.Client

var BotSetting = &Bot{}
var DbSetting = &Mysql{}

var Db *gorm.DB

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
	err = Cfg.Section("mysql").MapTo(&DbSetting)
	if err != nil {
		fmt.Println("Cfg.MapTo BotSetting err: ", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     RedisSetting.Host,
		Password: RedisSetting.Password,
		DB:       RedisSetting.DB})

	Db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DbSetting.User,
		DbSetting.Password,
		DbSetting.Host,
		DbSetting.Port,
		DbSetting.Dbname))
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	} else {
		log.Println("数据库连接成功")
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}

	Db.SingularTable(true)
	//Db.LogMode(true)
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)
}
