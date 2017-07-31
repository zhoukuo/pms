package main

import (
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter" //第三方库，用来处理HTTP的路由
	"github.com/zhoukuo/gosql"            //第三方库，用来处理SQL的操作
	"log"
	"net/http"
	"os"
	//以下以_开头的包均为业务逻辑处理相关的包
	_ "pms/business/event"
	_ "pms/business/harddelivery"
	_ "pms/business/hardstore"
	_ "pms/business/hardware"
	_ "pms/business/loan"
	_ "pms/business/project"
	_ "pms/business/returns"
	_ "pms/business/softdelivery"
	_ "pms/business/softstore"
	_ "pms/business/software"
	_ "pms/business/stock"
	_ "pms/business/user"
	"pms/config" //配置文件,包括所有环境变量,通用的数据结构定义等
)

func main() {
	//设置命令行参数
	port := flag.String("port", config.PORT, "http listen port")
	db := flag.String("db", config.DSN, "databse file")
	//从命令行中获取参数值
	flag.Parse()

	//验证数据库文件是否存在
	_, err := os.Stat(*db)
	if err != nil {
		msg := fmt.Sprintf("database(%s) not exist!", *db)
		panic(msg)
	}

	//连接数据库
	config.DB, _ = gosql.Open(config.DDN, *db)
	defer gosql.Close(config.DB)
	log.Printf("Initialize Database: DDN=%s, DSN=%s\n", config.DDN, *db)

	//初始化路由，每个业务的路由在各自的模块中初始化
	config.Router.GET("/", Index)

	//启动HTTP服务
	log.Println("Serving HTTP on 0.0.0.0 port", *port)
	http.ListenAndServe(":"+*port, config.Router)

}

//主页路由处理函数
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to PMS!")
}
