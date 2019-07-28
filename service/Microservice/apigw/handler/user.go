package handler

import (
	"FileStore-Server/service/Microservice/account/proto"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

var (
	userCli proto.UserService
)

func init()  {
	// 修改consul地址
	reg := consul.NewRegistry(func(op *registry.Options){
		op.Addrs = []string{
			"47.95.253.230:8500",
		}
	})

	service := micro.NewService(micro.Registry(reg), micro.Name("go.micro.api.user"))
	service.Init()

	// 初始化一个rpcClient
	userCli = proto.NewUserService("go.micro.service.user",service.Client())
}

//SignupHandler: 返回注册页面
func SignupHandler(c *gin.Context)  {
	c.Redirect(http.StatusFound,"/static/view/signup.html")
}

//DoSignupHandler: 处理用户注册
func DoSignupHandler(c *gin.Context)  {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")


	resp,err := userCli.Signup(context.TODO(),&proto.ReqSignup{
		Username:username,
		Password:passwd,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"code":resp.Code,
		"msg":resp.Message,
	})
}