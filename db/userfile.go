package db

import (
	mydb "FileStore-Server/db/mysql"
	"fmt"
	"time"
)

//UserFile: 用户文件表结构体
type UserFile struct {
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}

//OnUserFileUploadFinished: 更新用户文件表
func OnUserFileUploadFinished(username,filehash,filename string,filesize int64) bool {
	stmt,err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
			"`file_size`,`upload_at`) values (?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_,err = stmt.Exec(username,filehash,filename,filesize,time.Now())
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}