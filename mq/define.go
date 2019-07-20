package mq


// TransferData: 将要写到rabbitmq的数据的结构体
type TransferData struct {
	FileHash      string
	CurLocation   string	//当前存储位置
	DestLocation  string	//目标存储位置
}