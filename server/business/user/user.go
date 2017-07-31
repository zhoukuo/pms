package user

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/zhoukuo/gosql"
	"io/ioutil"
	"log"
	"net/http"
	"pms/config"
)

type User struct {
	Username    string
	Password    string
	Description string
}

func init() {
	// init router.
	config.Router.POST("/users/", AddUser)
	config.Router.POST("/verifyusers/", UserLogin)
	config.Router.GET("/users/", GetUserList)
	config.Router.OPTIONS("/user/", UserOptions)
	log.Println("Initialize user Runter ... ")
}

func UserOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func AddUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
	
	body, _ := ioutil.ReadAll(r.Body)
	var user User
	// all request must be json format
	json.Unmarshal(body, &user)
	if user.Username == "" || user.Password == ""{
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: 0, StatusText: "用户名和密码不能为空！"}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))
		return
	}

	sqlStatement := fmt.Sprintf("INSERT INTO tbl_user(username,password)VALUES('%s','%s')",user.Username,user.Password)
	log.Println(sqlStatement)

	id, _ := gosql.Insert(config.DB, sqlStatement)
	log.Println("returned rowid=", id)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "ok"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func UserLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
	
	body, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(body, &user)
	if user.Username == "" || user.Password == ""{
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: 0, StatusText: "用户名和密码不能为空！"}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))
		return
	}

	sqlStatement := fmt.Sprintf("SELECT username FROM tbl_user WHERE username='%s' and password='%s'",user.Username,user.Password)
	log.Println(sqlStatement)
	rows, _ := gosql.Query(config.DB, sqlStatement)
	responseText, _ := json.Marshal(rows)
	
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "ok"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func GetUserList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
	
	sqlStatement := fmt.Sprintf("SELECT username FROM tbl_user")
	rows, _ := gosql.Query(config.DB, sqlStatement)
	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "ok"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}
