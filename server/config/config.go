package config

import (
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

//动态(业务日志)操作类型的数据结构
type ActionStruct struct {
	CREATE string
	UPDATE string
}

//业务类型/设备状态的数据结构
type BusinessStruct struct {
	STORE    string
	DELIVERY string
	LOAN     string
	RETURN   string
	REJECT   string
	REPAIR   string
	FIXED    string
	INVALID  string
	DISUSE   string
}

//返回到web页面的状态码的数据结构
type StatusCodeStruct struct {
	COLUMN_NOT_SPECIFIED int
	UPDATE_NOT_ALLOW     int
	INVALID_NOT_ALLOW    int
	ID_ALREADY_EXIST     int
	ID_NOT_EXIST         int
	MODEL_EXIST          int
	VERSION_EXIST        int
	PLATFORM_EXIST       int
	PROJECT_EXIST        int
	STATUS_UNKNOWN       int
}

//返回到web页面的错误信息的数据结构
type StatusTextStruct struct {
	COLUMN_NOT_SPECIFIED string
	UPDATE_NOT_ALLOW     string
	INVALID_NOT_ALLOW    string
	ID_ALREADY_EXIST     string
	ID_NOT_EXIST         string
	MODEL_EXIST          string
	VERSION_EXIST        string
	PLATFORM_EXIST       string
	PROJECT_EXIST        string
	STATUS_UNKNOWN       string
}

//返回到web页面的数据结构
type ReturnStruct struct {
	ReadyStatus  int
	ResponseText string
	Status       int
	StatusText   string
}

var DB *sql.DB
var Router *httprouter.Router

//服务端环境配置
var DDN string = "sqlite3"
var DSN string = "pms.db"
var PORT string = "8088"

var OPERATOR string = "Admin"

//动态(业务日志)操作类型
var Action = ActionStruct{
	CREATE: "创建",
	UPDATE: "修改",
}

//业务类型/设备状态
var Business = BusinessStruct{
	STORE:    "入库",
	DELIVERY: "出库",
	LOAN:     "借出",
	RETURN:   "归还",
	REJECT:   "退货",
	REPAIR:   "维修",
	INVALID:  "作废",
	DISUSE:   "报废",
	FIXED:    "修复",
}

//返回到web页面的状态码
var StatusCode = StatusCodeStruct{
	COLUMN_NOT_SPECIFIED: 601,
	UPDATE_NOT_ALLOW:     602,
	INVALID_NOT_ALLOW:    603,
	ID_ALREADY_EXIST:     604,
	ID_NOT_EXIST:         605,
	MODEL_EXIST:          606,
	VERSION_EXIST:        607,
	PLATFORM_EXIST:       608,
	PROJECT_EXIST:        609,
	STATUS_UNKNOWN:       610,
}

//返回到web页面的错误信息
var StatusText = StatusTextStruct{
	COLUMN_NOT_SPECIFIED: "更新失败，未指定任何字段！",
	UPDATE_NOT_ALLOW:     "更新状态失败，当前状态不允许更新!",
	INVALID_NOT_ALLOW:    "作废失败，操作对象状态不允许作废！",
	ID_ALREADY_EXIST:     "硬件已存在！",
	ID_NOT_EXIST:         "硬件不存在！",
	MODEL_EXIST:          "软件登记失败，型号已存在！",
	VERSION_EXIST:        "软件入库失败，版本已存在！",
	PLATFORM_EXIST:       "硬件登记失败，平台已存在！",
	PROJECT_EXIST:        "项目登记失败，编号已存在！",
	STATUS_UNKNOWN:       "状态未知，请核对状态！",
}

func init() {
	Router = httprouter.New()
}

//允许任何跨域访问的配置
func SetAccessControlAllowOrigin(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Accept")
		//指定的时间内（单位是秒，以上设置的是 20 天）不需要再进行这种“事前检查”
		w.Header().Set("Access-Controll-Max-Age", "1728000")
		w.Header().Set("content-type", "application/json;charset=UTF-8")
	}
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func Guid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}
