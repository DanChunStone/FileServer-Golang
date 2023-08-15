package GinHandler

import (
	"FileStore-Server/common"
	"FileStore-Server/config"
	dblayer "FileStore-Server/db"
	nativeHandler "FileStore-Server/handler"
	"FileStore-Server/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// SignupHandler: 返回注册页面
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// DoSignupHandler: 处理用户注册
func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	if len(username) < 3 || len(passwd) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Invalid Parameter",
			"code": common.StatusParamInvalid,
		})
		return
	}

	//对密码进行加盐及取Sha1值加密
	encPasswd := util.Sha1([]byte(passwd + config.PwdSalt))
	ok := dblayer.UserSignUp(username, encPasswd)
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Signup succeeded",
			"code": common.StatusOK,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Signup failed",
			"code": common.StatusRegisterFailed,
		})
	}
}

// SignInHandler: 登录接口
func SignInHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

func PageNotFound(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/404.html") //这块不可以随便写，必须是302，否则就会报错
}

// DoSignInHandler: 处理登录请求
func DoSignInHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	//1.校验用户名与密码
	encPasswd := util.Sha1([]byte(password + config.PwdSalt))
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Name or password wrong",
			"code": common.StatusLoginFailed,
		})
		return
	}

	//2.生成用户访问凭证Token
	token := nativeHandler.GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "update token failed",
			"code": common.StatusLoginFailed,
		})
		return
	}

	//3.封装凭证与响应信息给客户端
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

// UserInfoHandler: 查询用户信息
func UserInfoHandler(c *gin.Context) {
	//1.解析请求参数
	username := c.Request.FormValue("username")

	//3.查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg":  "查询用户信息失败",
			"code": common.StatusUserNotExists,
		})
		return
	}

	//4.组装并响应用户数据
	resp := util.RespMsg{
		Code: int(common.StatusOK),
		Msg:  "查询用户信息成功",
		Data: user,
	}
	c.JSON(http.StatusOK, resp)
}
