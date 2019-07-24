package main

import (
	"FileStore-Server/service/Microservice/account/handler"
	"FileStore-Server/service/Microservice/account/proto"
	"github.com/micro/go-micro"
	"log"
)

func main()  {
	// 创建一个service
	service := micro.NewService(micro.Name("go.micro.service.user"))
	service.Init()

	proto.RegisterUserServiceHandler(service.Server(),new(handler.User))

	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
