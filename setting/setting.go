package setting

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"io"
	"log"
	"os"
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

type Plugin struct {
	Live    bool
	Danmaku bool
}

type Runtime struct {
	ImageMaxSize    int
	ImagePrefixUrl  string
	ImageSavePath   string
	RuntimeRootPath string
	ImageAllowExts  []string
}

// 初始化redis
var RedisSetting = &Redis{}
var RedisClient *redis.Client

// 初始化bot mysql
var BotSetting = &Bot{}
var DbSetting = &Mysql{}

// 全局mysql实例
var Db *gorm.DB

// 全局logrus实例
var Log = logrus.New()

// plugin设置
var PluginSetting = &Plugin{}

var RuntimeSetting = &Runtime{}

func SetUp() {
	// 初始化log
	Log.SetReportCaller(true)
	Log.SetFormatter(&logrus.TextFormatter{})
	today := time.Now().Format("2006102")
	file, err := os.OpenFile("runtime/logs/"+today+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("文件打开/创建失败")
	}
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	fileAndStd := io.MultiWriter(writers...)
	Log.SetOutput(fileAndStd)
	Log.SetLevel(logrus.InfoLevel)

	// 读取配置文件
	Cfg, err := ini.Load("conf/config.ini")
	if err != nil {
		Log.Fatal("Fail to parse 'conf/config.ini': ", err)
	}
	err = Cfg.Section("redis").MapTo(&RedisSetting)
	if err != nil {
		Log.Warning("Cfg.MapTo RedisSetting err: ", err)
	}
	err = Cfg.Section("bot").MapTo(&BotSetting)
	if err != nil {
		Log.Fatal("Cfg.MapTo BotSetting err: ", err)
	}
	err = Cfg.Section("mysql").MapTo(&DbSetting)
	if err != nil {
		Log.Warning("Cfg.MapTo BotSetting err: ", err)
	}
	err = Cfg.Section("plugin").MapTo(&PluginSetting)
	if err != nil {
		Log.Warning("Cfg.MapTo PluginSetting err: ", err)
	}
	err = Cfg.Section("runtime").MapTo(&RuntimeSetting)
	if err != nil {
		Log.Warning("Cfg.MapTo RuntimeSetting err: ", err)
	}

	// 初始化redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     RedisSetting.Host,
		Password: RedisSetting.Password,
		DB:       RedisSetting.DB})

	// 初始化mysql
	Db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DbSetting.User,
		DbSetting.Password,
		DbSetting.Host,
		DbSetting.Port,
		DbSetting.Dbname))
	if err != nil {
		Log.Error("models.Setup err: ", err)
	} else {
		Log.Info("数据库连接成功")
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}

	Db.SingularTable(true)
	Db.LogMode(true)
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)
}
