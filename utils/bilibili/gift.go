package bilibili

func SaveMsg() {
	mq := NewRabbitMQ("danmaku_test", "amq.direct", "")
	mq.ConsumeSimple()
}
