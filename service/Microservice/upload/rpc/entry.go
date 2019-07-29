package rpc

import (
	"FileStore-Server/service/Microservice/upload/config"
	uploadProto "FileStore-Server/service/Microservice/upload/proto"
	"context"
)

type Upload struct {} 

// 获取上传入口地址
func (u *Upload) UploadEntry(ctx context.Context,req *uploadProto.ReqEntry,resp *uploadProto.RespEntry) error {
	resp.Entry = config.UploadEntry
	return nil
}