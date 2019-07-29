package route

import (
	"FileStore-Server/service/Microservice/upload/api"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	// gin framework, 包括Logger, Recovery
	router := gin.Default()

	// 处理静态资源
	router.Static("/static/","./static")

	//// 加入中间件，用于验证token的拦截器
	//router.Use(handler.HTTPInterceptor())

	// 上传文件
	router.POST("/file/upload",api.DoUploadHandler)
	// 支持跨域
	router.OPTIONS("/file/upload", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 支持所有来源
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS") // 支持所有http方法
		c.Status(204) // 告诉前端请求成功，但body为空，页面不用刷新
	})

	return router
}