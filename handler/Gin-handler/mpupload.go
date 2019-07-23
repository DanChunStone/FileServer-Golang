package GinHandler

import (
	"fmt"
	"FileStore-Server/util"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	nativeHandler "FileStore-Server/handler"
	rPool "filestore-server/cache/redis"
	dblayer "FileStore-Server/db"
)

//InitialMultipartUploadHandler: 初始化分块上传
func InitialMultipartUploadHandler(c *gin.Context) {
	// 解析用户参数
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK,gin.H{
			"msg":"解析参数失败",
			"code":-1,
		})
		return
	}

	// 获得redis的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 生成分块上传的初始化信息
	upInfo := nativeHandler.MultipartUploadInfo{
		FileHash:filehash,
		FileSize:filesize,
		UploadID:username+fmt.Sprintf("%x",time.Now().UnixNano()),
		ChunkSize:5*1024*1024,
		ChunkCount:int(math.Ceil(float64(filesize)/(5*1024*1024))),
	}

	// 将初始化信息写入redis
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	// 将响应初始化数据返回到客户端
	c.Data(http.StatusOK,"application/json",util.NewRespMsg(0, "OK", upInfo).JSONBytes())
}

//UploadPartHandler: 上传文件分块
func UploadPartHandler(c *gin.Context) {
	// 解析用户请求参数
	// username := c.Request.FormValue("username")
	uploadID := c.Request.FormValue("uploadid")
	chunkIndex := c.Request.FormValue("index")

	// 获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 获得文件句柄，用于存储分块内容
	fpath := "tempFiles/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)

	fd, err := os.Create(fpath)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK,gin.H{
			"msg":"Upload part failed",
			"code":-1,
		})
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 返回处理结果到客户端
	c.Data(http.StatusOK,"application/json",util.NewRespMsg(0, "OK", nil).JSONBytes())
}

//CompleteUploadHandler: 通知上传合并
func CompleteUploadHandler(c *gin.Context) {
	// 解析请求参数
	upid := c.Request.FormValue("uploadid")
	username := c.Request.FormValue("username")
	filehash := c.Request.FormValue("filehash")
	filesize := c.Request.FormValue("filesize")
	filename := c.Request.FormValue("filename")

	// 获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 通过uploadid查询redis并判断是否所有分块上传完成
	data,err := redis.Values(rConn.Do("HGETALL","MP_"+upid))
	if err != nil{
		log.Println(err.Error())
		c.JSON(http.StatusOK,gin.H{
			"msg":"complete upload failed",
			"code":-1,
		})
		return
	}

	totalCount := 0 //所有块的数量
	chunkCount := 0	//已完成的块的数量
	for i:=0;i<len(data);i+=2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount,_ = strconv.Atoi(v)
		}else if strings.HasPrefix(k,"chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount == chunkCount {
		c.JSON(http.StatusOK,gin.H{
			"msg":  "Complete failed, invalid request",
			"code": -2,
		})
		return
	}

	// TODO：合并分块

	// 更新唯一文件表及用户文件表
	fsize,_ := strconv.Atoi(filesize)
	dblayer.OnFileUploadFinished(filehash,filename,int64(fsize),"")
	dblayer.OnUserFileUploadFinished(username,filehash,filename,int64(fsize))

	// 响应处理结果
	c.Data(http.StatusOK,"application/json",util.NewRespMsg(0, "OK", nil).JSONBytes())
}