package rpc

import (
	"context"
	cfg "FileStore-Server/service/Microservice/download/config"
	dlProto "FileStore-Server/service/Microservice/download/proto"
)

// Dwonload :download结构体
type Download struct{}

// DownloadEntry : 获取下载入口
func (u *Download) DownloadEntry(
	ctx context.Context,
	req *dlProto.ReqEntry,
	res *dlProto.RespEntry) error {

	res.Entry = cfg.DownloadEntry
	return nil
}
