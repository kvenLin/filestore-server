package db

import (
	"fmt"
	"imooc.com/course/filestore-server/db/mysql"
)

//通过用户名和密码完成用户表的注册操作
func UserSignUp(username string, password string) bool {
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_user (`user_name`, `user_pwd`) values(?, ?)")
	if err != nil {
		fmt.Println("Failed to insert , err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert , err:" + err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

//用户登录
func UserSignIn(username string, encpwd string) bool {
	stmt, err := mysql.DBConn().Prepare("select * from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("Username not found: " + username)
		return false
	}
	//todo pare rows
	return true
}

//刷新用户登录的token
func UpdateToken(username string, token string) bool {
	stmt, err := mysql.DBConn().Prepare("replace into tbl_user_token (`user_name`, `user_token`) values(?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

//用户信息查询
func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mysql.DBConn().Prepare(
		"select user_name, signup_at from tbl_user where user_name = ? limit 1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
