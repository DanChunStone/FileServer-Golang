package mq

import (
	"github.com/streadway/amqp"
	"log"
)

// Publish: 发布消息
func Publish(exchange, routingKey string, msg []byte) bool {
	if !initChannel() {
		return false
	}

	// 通过channel发布消息
	err := channel.Publish(exchange, routingKey,
		false,	// 如果没有对应的queue, 就会丢弃这条消息
		false,		// 最新rabbitmq中已不起作用
		amqp.Publishing{
			ContentType: "text/plain",	//明文
			Body: msg,
		})
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}
