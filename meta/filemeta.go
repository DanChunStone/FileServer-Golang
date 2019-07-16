package meta

import (
	mydb "FileStore-Server/db"
)

// FileMeta: 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

// fileMetas: 以文件哈希值为key，文件元信息结构体为value
var fileMetas map[string]FileMeta

func init()  {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta: 新增/更新文件元信息
func UpdateFileMeta(fMeta FileMeta)  {
	fileMetas[fMeta.FileSha1] = fMeta
}

// UpdateFileMetaDB: 新增/更新元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool{
	return mydb.OnFileUploadFinished(fmeta.FileSha1,fmeta.FileName,fmeta.FileSize,fmeta.Location)
}

// UpdateFileMetaDB: 从mysql中获取文件元信息
func GetFileMetaDB(fileSha1 string) (*FileMeta,error){
	tfile,err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return &FileMeta{},err
	}

	fmeta := FileMeta{
		FileSha1:tfile.FileHash,
		FileName:tfile.FileName.String,
		FileSize:tfile.FileSize.Int64,
		Location:tfile.FileAddr.String,
	}

	return &fmeta,nil
}

// GetFileMeta: 通过sha1获取文件元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// RemoveFileMeta: 删除元信息
func RemoveFileMeta(fileSha1 string)  {
	delete(fileMetas,fileSha1)
}
