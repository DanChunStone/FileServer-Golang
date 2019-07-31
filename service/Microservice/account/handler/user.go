package handler

import (
	"FileStore-Server/common"
	"FileStore-Server/config"
	dblayer "FileStore-Server/db"
	proto "FileStore-Server/service/Microservice/account/proto"
	"FileStore-Server/util"
	"context"
	"fmt"
	"time"
)

type User struct {}

// GenToken : 生成token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

//Signup: 用户注册
func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, resp *proto.RespSignup) error {
	username := req.Username
	passwd := req.Password

	if len(username) < 3 || len(passwd) < 5 {
		resp.Code = common.StatusParamInvalid
		resp.Message = "注册参数无效"
		return nil
	}

	//对密码进行加盐及取Sha1值加密
	encPasswd := util.Sha1([]byte(passwd+config.PwdSalt))
	ok := dblayer.UserSignUp(username,encPasswd)
	if ok {
		resp.Code = common.StatusOK
		resp.Message = "注册成功"
	}else {
		resp.Code = common.StatusRegisterFailed
		resp.Message = "注册失败"
	}
	return nil
}

// 用户登录
func (u *User) Signin(context.Context, *proto.ReqSignin, *proto.RespSignin) error {

}

// 获取用户信息
func (u *User) UserInfo(context.Context, *proto.ReqUserInfo, *proto.RespUserInfo) error {

}