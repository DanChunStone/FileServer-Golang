package GinHandler

import (
	"FileStore-Server/common"
	"FileStore-Server/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	nativeHandler "FileStore-Server/handler"
)

//HTTPInterceptor: http请求拦截器
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		log.Println("\nusername:"+username+"\ntoken:"+token)

		// 验证token是否有效
		if len(username)<3 || !nativeHandler.IsTokenValid(token) {
			// 如果token校验失败，则跳转到直接返回失败提示
			c.Abort()
			c.JSON(http.StatusOK,util.NewRespMsg(int(common.StatusTokenInvalid), "token无效", nil, ))
			return
		}

		// 如果校验通过，则将请求转到下一个handler处理
		c.Next()
	}
}