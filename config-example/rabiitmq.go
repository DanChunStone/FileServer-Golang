package config

const (
	// AsyncTransferEnable : 是否开启文件异步转移(默认同步)
	AsyncTransferEnable = true
	// TransExchangeName : 用于文件transfer的交换机
	TransExchangeName = "uploadserver.name"
	// TransOSSQueueName : oss转移队列名
	TransOSSQueueName = "uploadserver.name.oss"
	// TransOSSErrQueueName : oss转移失败后写入另一个队列的队列名
	TransOSSErrQueueName = "uploadserver.name.oss.err"
	// TransOSSRoutingKey : routingkey
	TransOSSRoutingKey = "name"
)

var (
	// RabbitURL : rabbitmq服务的入口url
	RabbitURL = "amqp://guest:guest@localhost:5672/"
)