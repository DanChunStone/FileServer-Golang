package GinHandler

import (
	"FileStore-Server/config"
	"FileStore-Server/db"
	"FileStore-Server/meta"
	"FileStore-Server/mq"
	"FileStore-Server/store/oss"
	"FileStore-Server/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//UploadHandler: 文件上传
func UploadHandler(c *gin.Context) {
	//返回上传文件的html页面
	c.Redirect(http.StatusFound,"/static/view/index.html")
}

//DoUploadHandler: 处理文件上传
func DoUploadHandler(c *gin.Context) {
	errCode := 0
	defer func() {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
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
		//重定向至主页
		c.Redirect(http.StatusFound,"/static/view/home.html")
	}else {
		errCode = -4
	}

	//log
	log.Printf("Your file's meta is: %s\n",fileMeta)
}

//GetFileMetaHandler: 通过文件hash值，获取文件元信息
func GetFileMetaHandler(c *gin.Context)  {
	//获取hash值，并通过其查询文件元信息
	filehash := c.Request.Form["filehash"][0]
	fMeta,err := meta.GetFileMetaDB(filehash)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK,gin.H{
			"msg":"文件元信息不存在",
			"code":-1,
		})
		return
	}

	//json格式化meta实例
	if fMeta != nil {
		data, err := json.Marshal(fMeta)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK,gin.H{
				"msg":"文件元信息格式化失败",
				"code":-2,
			})
			return
		}
		c.Data(http.StatusOK,"application/json",data)
	} else {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"no such file",
		})
	}
}

// FileQueryHandler: 查询批量的文件元信息
func FileQueryHandler(c *gin.Context) {
	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")

	userFiles, err := db.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK,gin.H{
			"msg":"文件元信息不存在",
			"code":-1,
		})
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK,gin.H{
			"msg":"文件元信息格式化失败",
			"code":-2,
		})
		return
	}
	c.Data(http.StatusOK,"application/json",data)
}

//DownloadHandler: 根据文件哈希值下载文件
func DownloadHandler(c *gin.Context)  {
	//获取文件hash值，并获取元信息
	fileSha1 := c.Request.FormValue("filehash")
	fm,err := meta.GetFileMetaDB(fileSha1)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK,gin.H{
			"msg":"文件元信息不存在",
			"code":-1,
		})
		return
	}

	c.FileAttachment(fm.Location,fm.FileName)
}

//FileMetaUpdateHandler: 更新文件元信息
func FileMetaUpdateHandler(c *gin.Context) {
	opType := c.Request.FormValue("op")
	fileSha1 := c.Request.FormValue("filehash")
	newFileName := c.Request.FormValue("filename")

	if opType != "0" {
		c.JSON(http.StatusOK,gin.H{
			"msg":"参数错误",
			"code":-1,
		})
		return
	}

	// 更新文件元信息
	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data,err := json.Marshal(curFileMeta)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK,gin.H{
			"msg":"文件元信息格式化失败",
			"code":-2,
		})
		return
	}
	c.Data(http.StatusOK,"application/json",data)
}

//FileDeleteHandler: 删除文件
func FileDeleteHandler(c *gin.Context) {
	fileSha1 := c.Request.FormValue("filehash")

	//删除本地文件
	fMeta := meta.GetFileMeta(fileSha1)
	err := os.Remove(fMeta.Location)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK,gin.H{
			"msg":"服务器删除文件失败",
			"code":-1,
		})
		return
	}

	//删除元信息
	meta.RemoveFileMeta(fileSha1)

	c.Writer.WriteHeader(http.StatusOK)
}

//TryFastUploadHandler: 尝试进行秒传
func TryFastUploadHandler(c *gin.Context) {
	//解析请求参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filename := c.Request.FormValue("filename")
	filesize,_ := strconv.Atoi(c.Request.FormValue("filesize"))

	//从文件表中查询相同hash的文件记录
	fileMeta,err := meta.GetFileMetaDB(filehash)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK,gin.H{
			"msg":"查询文件记录失败",
			"code":-1,
		})
		return
	}

	if fileMeta == nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"没有找到相同文件，请尝试普通接口",
			"code":-2,
		})
		return
	}

	ok := db.OnUserFileUploadFinished(username,filehash,filename,int64(filesize))
	if ok {
		c.JSON(http.StatusOK,gin.H{
			"msg":"秒传成功",
			"code":0,
		})
		return
	}else {
		c.JSON(http.StatusOK,gin.H{
			"msg":"秒传失败，请稍后再试或尝试普通接口",
			"code":-3,
		})
		return
	}
}

//DownloadURLHandler: 获取oss文件下载URL
func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")

	// 从文件表查找记录
	row, _ := db.GetFileMeta(filehash)

	// 判断是否存储在oss，并获取url
	signedURL := oss.DownloadURL(row.FileAddr.String)

	c.Writer.Write([]byte(signedURL))
}