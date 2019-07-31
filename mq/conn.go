package mq

import (
	"FileStore-Server/config"
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 如果异常关闭，会接收通知
var notifyClose chan *amqp.Error

// UpdateRabbitHost : 更新mq host
func UpdateRabbitHost(host string) {
	config.RabbitURL = host
}

func Init() {
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