package bilibili

import (
	"errors"
	"fmt"
	"go-bot/setting"
	"log"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
	"github.com/tidwall/gjson"
)

// MQURL url格式  amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://go_bot:0054444944@123.60.53.237:5672/go-bot"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	// 连接信息
	Mqurl string
}

// NewRabbitMQ 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
	var err error
	// 创建rabbitmq连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误！")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败！")
	return rabbitmq
}

// Destory 断开channel和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// NewRabbitMQSimple 简单模式step1： 1.创建简单模式下的rabbitmq实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "amq.direct", "")
}

// PublishSimple 简单模式step2： 2.简单模式下生产
func (r *RabbitMQ) PublishSimple(message string) {
	// 1. 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true,根据exchange类型和routekey规则，如果无法找到符合条件的队列，则会把发送的消息返回给发送者
		false,
		// 如果为true,当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息发还给发送者
		false,
		amqp.Publishing{ContentType: "application/json", Body: []byte(message)},
	)
}

type GiftLog struct {
	Nickname string
	GiftName string
	Price    float64
	Num      float64
	gorm.Model
}

type MqLog struct {
	LogData string
}

// ConsumeSimple 简单模式step3： 3.简单模式下消费
func (r *RabbitMQ) ConsumeSimple() {
	fmt.Println("开始接收消息")
	var msg struct {
		cmd  string
		data map[string]interface{}
	}
	// 1. 接收消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果为true，表示不能将同一个conn中的消息发送给这个conn中的消费者
		false,
		// 队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	// 2. 启用协程处理消息
	go func() {
		for d := range msgs {
			// 实现我们要处理的逻辑函数
			dd := fmt.Sprintf("%s", d.Body)
			test := gjson.Parse(dd).Value().(map[string]interface{})
			msg.cmd = test["cmd"].(string)
			msg.data = test["data"].(map[string]interface{})
			glog := MqLog{
				LogData: dd,
			}
			saveLog(glog)
			fmt.Println(reflect.TypeOf(msg.data["price"]))

			if msg.cmd == "GUARD_BUY" {
				data := GiftLog{
					Nickname: msg.data["username"].(string),
					GiftName: msg.data["gift_name"].(string),
					Num:      msg.data["num"].(float64),
					Price:    msg.data["price"].(float64),
				}
				err = gitfAdd(data)
			} else {
				data := GiftLog{
					Nickname: msg.data["uname"].(string),
					GiftName: msg.data["giftName"].(string),
					Num:      msg.data["num"].(float64),
					Price:    msg.data["price"].(float64),
				}
				err = gitfAdd(data)
			}

			// if giftExsit(data) {

			if err != nil {
				log.Println("[gift add err]:" + err.Error())
			}
			// }
		}
	}()
	log.Printf("[*] waiting for messages, to exit process CTRL+C")
	<-forever
}

func giftExsit(data GiftLog) bool {
	fmt.Println(data)
	if errors.Is(setting.Db.Where(&data).First(&GiftLog{}).Error, gorm.ErrRecordNotFound) {
		return true
	}
	return false
}

func gitfAdd(data GiftLog) error {
	return setting.Db.Create(&data).Error
}

func saveLog(data MqLog) error {
	return setting.Db.Create(&data).Error
}
