package db

import (
	"database/sql"
	"fmt"
	"imooc.com/course/filestore-server/db/mysql"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
	UploadAt string
}

//文件上传完成保存Meta
func OnFileUploadFinished(filehash string, filename string,
	filesize int64, fileaddr string) bool {
	stmt, err := mysql.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) values(?,?,?,?,1)",
	)
	if err != nil {
		fmt.Println("Failed to paper statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf < 0 {
			fmt.Printf("File with hash: %s has uploaded before ", filehash)
		}
		return true
	}
	return false
}

//从数据库中获取文件元数据
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"select file_sha1, file_name, file_size, file_addr from tbl_file where file_sha1 = ? and status = 1 limit 1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()
	tfile := TableFile{}
	//对查询结果进行赋值
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, nil
}
