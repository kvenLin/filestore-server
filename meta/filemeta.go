package meta

import "imooc.com/course/filestore-server/db"

//文件元数据结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

//初始化
func init() {
	fileMetas = make(map[string]FileMeta)

}

//新增或者更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

//通过sha1获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

//更新元信息到DB
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

//从mysql获取文件元信息
func GetFileMetaDB(filahash string) (FileMeta, error) {
	tfile, err := db.GetFileMeta(filahash)
	if err != nil {
		return FileMeta{}, err
	}
	fMeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fMeta, nil
}
