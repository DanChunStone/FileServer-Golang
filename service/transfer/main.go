package main

import (
	"FileStore-Server/config"
	"FileStore-Server/db"
	"FileStore-Server/mq"
	"FileStore-Server/store/oss"
	"bufio"
	"encoding/json"
	"log"
	"os"
)

//ProcessTransfer: 处理文件转移
func ProcessTransfer(msg []byte) bool {
	log.Println(string(msg))

	// 解析msg
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 根据msg临时存储位置，并创建文件句柄
	filed, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 通过文件句柄将文件内容转移到oss
	err = oss.Bucket().PutObject(pubData.DestLocation,bufio.NewReader(filed))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// 更新文件的存储路径到数据库文件表
	ok := db.UpdateFileLocation(pubData.FileHash,pubData.DestLocation)
	if !ok {
		return false
	}

	return true
}

func main()  {
	log.Println("开始监听转移任务队列")
	mq.StartConsume(config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer,
	)
}
