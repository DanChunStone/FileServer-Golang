package handler

import (
	"FileStore-Server/util"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	dblayer "FileStore-Server/db"
	"time"
)

const (
	//用于加密的盐值(自定义)
	pwdSalt = "*#890"
)

//SignupHandler : 处理用户注册请求
func SignupHandler(w http.ResponseWriter,r *http.Request) {
	//GET: 页面请求
	if r.Method == http.MethodGet {
		data,err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	//POST: 数据请求
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.Form.Get("username")
		passwd := r.Form.Get("password")

		if len(username) < 3 || len(passwd) < 5 {
			w.Write([]byte("Invalid Parameter"))
			return
		}

		//对密码进行加盐及取Sha1值加密
		encPasswd := util.Sha1([]byte(passwd+pwdSalt))
		ok := dblayer.UserSignUp(username,encPasswd)
		if ok {
			w.Write([]byte("SUCCESS"))
		}else {
			w.Write([]byte("FAILED"))
		}
	}
}

//SignInHandler: 登录接口
func SignInHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	//1.校验用户名与密码
	encPasswd := util.Sha1([]byte(password + pwdSalt))
	pwdChecked := dblayer.UserSignin(username,encPasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	//2.生成用户访问凭证Token
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username,token)
	if !upRes {
		fmt.Println("更新访问凭证失败")
		w.Write([]byte("FAILED"))
		return
	}

	//3.封装凭证与响应信息给客户端
	resp := util.RespMsg{
		Code:0,
		Msg:"OK",
		Data: struct {
			Location	string
			Username	string
			Token		string
		}{
			Location:"http://"+r.Host+"/static/view/home.html",
			Username:username,
			Token:token,
		},
	}
	w.Write(resp.JSONBytes())
}

//UserInfoHandler: 查询用户信息
func UserInfoHandler(w http.ResponseWriter,r *http.Request) {
	//1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")

	//以下功能已经放入拦截器中
	////2.验证Token是否有效
	//isValidToken := IsTokenValid(token)
	//if !isValidToken {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}

	//3.查询用户信息
	user,err := dblayer.GetUserInfo(username)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//4.组装并响应用户数据
	resp := util.RespMsg{
		Code:0,
		Msg:"OK",
		Data:user,
	}
	w.Write(resp.JSONBytes())
}

//GenToken: 生成40位token
func GenToken(username string) string {
	//40位token: md5(username + timestamp(时间戳) + token_salt) + timestamp[:8]
	ts := fmt.Sprint("%x",time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix + ts[:8]
}

//IsTokenValid: 验证token是否失效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		fmt.Println("token长度异常")
		return false
	}
	return true
}