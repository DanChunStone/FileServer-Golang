package mq

import (
	"log"
)

var done chan bool

//StartConsume: 开始监听消息队列，获取消息并调用回调函数
func StartConsume(queueName,channelName string,callback func(msg []byte) bool)  {
	// 通过channel.Consume获取消息通道
	msgs, err := channel.Consume(
		queueName,
		channelName,
		true,	// 自动应答
		false,	// 非唯一的消费者
		false,	// rabbitMQ只能设置为false
		false,	// noWait, false表示会阻塞直到有消息过来
		nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	done = make(chan bool)

	// 循环获取队列的消息
	go func() {
		for msg := range msgs {
			// 对每个消息，调用回调函数来处理消息
			processOK := callback(msg.Body)
			if !processOK {
				// TODO: 将任务写到另一个队列，用于异常情况的重试
			}
		}
	}()

	// 直到读取队列成功，都会阻塞
	<- done

	// 关闭通道
	channel.Close()
}

//StopConsume: 停止监听队列
func StopConsume() {
	done <- true
}