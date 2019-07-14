package db

import (
	mydb "FileStore-Server/db/mysql"
	"fmt"
)

//User: 用户表model
type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

//UserSignUp: 通过用户名和密码完成注册
func UserSignUp(username string,passwd string)bool {
	//预编译sql语句
	stmt,err := mydb.DBConn().Prepare(
			"insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("Failed to insert, err: "+err.Error())
		return false
	}

	//进行数据库插入
	ret,err := stmt.Exec(username,passwd)
	if err != nil {
		fmt.Println("Failed to insert, err: "+err.Error())
	}

	//判断注册结果并返回
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

//UserSignin: 验证登录
func UserSignin(username string,encpwd string) bool {
	//预编译sql语句
	stmt,err := mydb.DBConn().Prepare(
		"select * from tbl_user where user_name = ? limit 1")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//进行数据库查询，并处理错误
	rows,err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Printf("username:%s not found: \n",username)
		return false
	}

	//将查询结果转换为map数组
	pRows := mydb.ParseRows(rows)
	//判断查询结果长度以及密码是否匹配
	if len(pRows)>0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

//UpdateToken: 刷新用户登录token
func UpdateToken(username string,token string) bool {
	stmt,err := mydb.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values (?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_,err = stmt.Query(username,token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

//GetUserInfo: 查询用户信息
func GetUserInfo(username string) (User,error) {
	user := User{}

	stmt,err := mydb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	defer stmt.Close()
	if err != nil {
		return user,err
	}

	//执行查询
	err = stmt.QueryRow(username).Scan(&user.Username,&user.SignupAt)
	if err != nil {
		return user,err
	}

	return user,nil
}