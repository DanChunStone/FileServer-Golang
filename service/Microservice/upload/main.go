package main

import (
	cfg "FileStore-Server/service/Microservice/upload/config"
	upProto "FileStore-Server/service/Microservice/upload/proto"
	"FileStore-Server/service/Microservice/upload/route"
	upRpc "FileStore-Server/service/Microservice/upload/rpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"log"
	"time"
)

func startRPCService() {
	// 修改consul地址
	reg := consul.NewRegistry(func(op *registry.Options){
		op.Addrs = []string{
			"47.95.253.230:8500",
		}
	})

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.upload"),	// 服务名称
		micro.RegisterTTL(time.Second*10),			// TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5),		// 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
	)
	service.Init()

	upProto.RegisterUploadServiceHandler(service.Server(), new(upRpc.Upload))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

// 启动api服务
func startAPIService() {
	router := route.Router()
	router.Run(cfg.UploadServiceHost)
}

func main()  {
	go startRPCService()
	startAPIService()
}