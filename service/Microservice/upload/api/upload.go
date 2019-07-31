package api

import (
	"FileStore-Server/config"
	"FileStore-Server/db"
	"FileStore-Server/meta"
	"FileStore-Server/mq"
	dbcli "FileStore-Server/service/Microservice/dbproxy/client"
	"FileStore-Server/service/Microservice/dbproxy/orm"
	"FileStore-Server/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//DoUploadHandler: 处理文件上传
func DoUploadHandler(c *gin.Context) {
	errCode := 0
	defer func() {
		// 处理跨域请求
		c.Header("Access-Control-Allow-Origin", "*") // 支持所有来源
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS") // 支持所有http方法
		// 处理上传结果
		if errCode < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "上传失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": errCode,
				"msg":  "上传成功",
			})
		}
	}()

	//获取表单上传的文件，并打开
	file,head,err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Failed to get data, err: %s\n",err.Error())
		errCode = -1
		return
	}
	defer file.Close()

	//创建文件元信息实例
	fileMeta := meta.FileMeta{
		FileName:head.Filename,
		Location:"tempFiles/"+head.Filename,
		UploadAt:time.Now().Format("2006-01-02 15:04:05"),
	}

	//创建本地文件
	localFile,err := os.Create(fileMeta.Location)
	if err != nil {
		log.Printf("Failed to create file, err: %s\n",err.Error())
		errCode = -2
		return
	}
	defer localFile.Close()

	//复制文件信息到本地文件
	fileMeta.FileSize,err = io.Copy(localFile,file)
	if err != nil {
		log.Printf("Failed to save data into file, err: %s\n",err.Error())
		errCode = -3
		return
	}

	//计算文件哈希值
	localFile.Seek(0,0)
	fileMeta.FileSha1 = util.FileSha1(localFile)

	// 游标重新回到文件头部
	localFile.Seek(0, 0)

	// 将文件写入阿里云OSS
	ossPath := "oss/" + fileMeta.FileSha1 + fileMeta.FileName

	// 将转移任务添加到rabbitmq队列中
	data := mq.TransferData{
		FileHash:fileMeta.FileSha1,
		CurLocation:fileMeta.Location,
		DestLocation:ossPath,
	}
	pubData, _ := json.Marshal(data)
	pubSuc := mq.Publish(
		config.TransExchangeName,
		config.TransOSSRoutingKey,
		pubData,
	)
	if !pubSuc {
		// TODO: 当前发送转移信息失败，稍后重试
	}

	//将文件元信息添加到mysql中
	_ = meta.UpdateFileMetaDB(fileMeta)

	//更新用户文件表记录
	username := c.Request.FormValue("username")
	ok := db.OnUserFileUploadFinished(username,fileMeta.FileSha1,fileMeta.FileName,fileMeta.FileSize)
	if ok {
		errCode = 0
	}else {
		errCode = -4
	}
}

// TryFastUploadHandler : 尝试秒传接口
func TryFastUploadHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	//filesize, _ := strconv.Atoi(c.Request.FormValue("filesize"))

	// 2. 从文件表中查询相同hash的文件记录
	fileMetaResp, err := dbcli.GetFileMeta(filehash)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 3. 查不到记录则返回秒传失败
	if !fileMetaResp.Suc {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}

	// 4. 上传过则将文件信息写入用户文件表， 返回成功
	fmeta := dbcli.TableFileToFileMeta(fileMetaResp.Data.(orm.TableFile))
	fmeta.FileName = filename
	upRes, err := dbcli.OnUserFileUploadFinished(username, fmeta)
	if err == nil && upRes.Suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		c.Data(http.StatusOK, "application/json", resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
	return
}