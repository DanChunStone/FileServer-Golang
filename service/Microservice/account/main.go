package main

import (
	"FileStore-Server/service/Microservice/account/handler"
	"FileStore-Server/service/Microservice/account/proto"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"log"
	"time"
)

func main()  {
	// 修改consul地址
	reg := consul.NewRegistry(func(op *registry.Options){
		op.Addrs = []string{
			"47.95.253.230:8500",
		}
	})

	// 创建一个service
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.user"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		)
	service.Init()

	proto.RegisterUserServiceHandler(service.Server(),new(handler.User))

	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
