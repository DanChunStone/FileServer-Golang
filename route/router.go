package route

import (
	handler "FileStore-Server/handler/Gin-handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	// gin framework, 包括Logger, Recovery
	router := gin.Default()

	// 处理静态资源
	router.Static("/static/","./static")

	// 无需验证就能访问的接口
	router.GET("/user/signup",handler.SignupHandler)
	router.POST("/user/signup",handler.DoSignupHandler)
	router.GET("/user/signin",handler.SignInHandler)
	router.POST("/user/signin",handler.DoSignInHandler)
	router.GET("/", handler.SignInHandler)

	// 加入中间件，用于验证token的拦截器
	router.Use(handler.HTTPInterceptor())

	// Use 中间件之后的接口，都需要通过拦截器
	// 用户信息
	router.POST("/user/info",handler.UserInfoHandler)

	// 上传文件
	router.GET("/file/upload",handler.UploadHandler)
	router.POST("/file/upload",handler.DoUploadHandler)
	router.POST("/file/fastupload",handler.TryFastUploadHandler)
	// 查询文件
	router.POST("/file/meta",handler.GetFileMetaHandler)
	router.POST("/file/query",handler.FileQueryHandler)
	// 下载文件
	router.POST("/file/download",handler.DownloadHandler)
	router.POST("/file/downloadurl",handler.DownloadURLHandler)
	// 更新文件
	router.POST("/file/update",handler.FileMetaUpdateHandler)
	// 删除文件
	router.POST("/file/delete",handler.FileDeleteHandler)

	// 分块上传
	router.POST("/file/mpupload/init",handler.InitialMultipartUploadHandler)
	router.POST("/file/mpupload/uppart",handler.UploadPartHandler)
	router.POST("/file/mpupload/complete",handler.CompleteUploadHandler)

	return router
}