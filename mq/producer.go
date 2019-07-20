package mq

import (
	"FileStore-Server/config"
	"github.com/streadway/amqp"
	"log"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 如果异常关闭，会接收通知
var notifyClose chan *amqp.Error

func init() {
	// 是否开启异步转移功能，开启时才初始化rabbitMQ连接
	if !config.AsyncTransferEnable {
		return
	}
	if initChannel() {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel()
			}
		}
	}()
}

func initChannel() bool {
	// 判断channel是否已经创建
	if channel != nil {
		return true
	}

	// 获取rabbitmq的一个连接
	conn,err := amqp.Dial(config.RabbitURL)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 打开一个channel，用于消息的发布与接收
	channel,err = conn.Channel()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

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
