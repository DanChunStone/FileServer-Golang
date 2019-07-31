package main

import (
	"FileStore-Server/common"
	"FileStore-Server/service/Microservice/dbproxy/config"
	"log"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"

	dbConn "FileStore-Server/service/Microservice/dbproxy/conn"
	dbProxy "FileStore-Server/service/Microservice/dbproxy/proto"
	dbRpc "FileStore-Server/service/Microservice/dbproxy/rpc"
)

func startRpcService() {
	service := micro.NewService(
		micro.Name("go.micro.service.dbproxy"),	// 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),			// 声明超时时间, 避免consul不主动删掉已失去心跳的服务节点
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			// 检查是否指定dbhost
			dbhost := c.String("dbhost")
			if len(dbhost) > 0 {
				log.Println("custom db address: " + dbhost)
				config.UpdateDBHost(dbhost)
			}
		}),
	)

	// 初始化db connection
	dbConn.InitDBConn()

	dbProxy.RegisterDBProxyServiceHandler(service.Server(), new(dbRpc.DBProxy))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

func main() {
	startRpcService()
}

// res, err := mapper.FuncCall("/user/UserExist", []interface{}{"haha"}...)
// log.Printf("error: %+v\n", err)
// log.Printf("result: %+v\n", res[0].Interface())

// res, err = mapper.FuncCall("/user/UserExist", []interface{}{"admin"}...)
// log.Printf("error: %+v\n", err)
// log.Printf("result: %+v\n", res[0].Interface())
