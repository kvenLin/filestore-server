package handler

import (
	"fmt"
	"imooc.com/course/filestore-server/db"
	"imooc.com/course/filestore-server/util"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "#890"
)

//处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	//校验
	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invalid parameter!"))
		return
	}

	//对密码加密
	encPassword := util.Sha1([]byte(password + pwd_salt))
	suc := db.UserSignUp(username, encPassword)
	if suc {
		w.Write([]byte("SUCCESS!"))
	} else {
		w.Write([]byte("FAILED!"))
	}
}

//登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPassword := util.Sha1([]byte(password + pwd_salt))
	//1.校验用户名和密码
	pwdChecked := db.UserSignIn(username, encPassword)
	if !pwdChecked {
		w.Write([]byte("FAILED!"))
		return
	}
	//2.生成访问的凭证token
	token := GenToken(username)
	upRes := db.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED!"))
		return
	}
	//3.登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

//生成token
func GenToken(username string) string {
	//md5(username + timestamp + salt) + timestamp[:8]
	timestamp := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + timestamp + "_tokesalt"))
	return tokenPrefix + timestamp[:8]
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")

	//2.验证token是否有效
	//valid := IsTokenValid(username, token)
	//if !valid {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	//3.查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//4.组装并响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

//验证token是否有效
func IsTokenValid(username string, token string) bool {
	//1.判断token的时效性,是否过期

	//2.数据库表tbl_user_token查询username对应的token信息

	//3.直接对比token是否一致,一致返回true,否则返回false
	return true
}
