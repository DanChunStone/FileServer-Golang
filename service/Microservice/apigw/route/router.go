package route

import (
	"FileStore-Server/service/Microservice/apigw/handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	// 处理静态资源
	router.Static("/static/","./static")

	router.GET("/user/signup",handler.SignupHandler)
	router.POST("/user/signup",handler.DoSignupHandler)

	return router
}