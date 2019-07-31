package main

import (
	"fmt"
	"time"

	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"

	"FileStore-Server/common"
	dbproxy "FileStore-Server/service/Microservice/dbproxy/client"
	cfg "FileStore-Server/service/Microservice/download/config"
	dlProto "FileStore-Server/service/Microservice/download/proto"
	"FileStore-Server/service/Microservice/download/route"
	dlRpc "FileStore-Server/service/Microservice/download/rpc"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.download"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlRpc.Download))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startAPIService() {
	router := route.Router()
	router.Run(cfg.DownloadServiceHost)
}

func main() {
	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
