package harddelivery

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/zhoukuo/gosql"
	"io/ioutil"
	"log"
	"net/http"
	"pms/business/event"
	"pms/business/stock"
	"pms/config"
	"strings"
	"time"
)

type Harddeliverys struct {
	Order_id        string
	Date            string
	ID              string
	Price           string
	Project_id      string
	Soft_model      string
	Version         string
	Algorithm       string
	Period          string
	User            string
	Delivery_type   string
	Tracking_number string
	Description     string
	Platform        string
	Model           string
}

func init() {
	// init router.
	config.Router.POST("/harddeliverys/", HarddeliveryNew)
	config.Router.POST("/harddeliverysrejected/", HarddeliveryRejected)
	config.Router.GET("/harddeliverys/", HarddeliverysGet)
	config.Router.GET("/harddeliverys/:guid", HarddeliverysGetByGUID)
	config.Router.POST("/harddeliverys/:guid", HarddeliverysUpdateByGUID)
	config.Router.GET("/harddelivery_order_list/", HarddeliveryOrderList)
	config.Router.GET("/harddelivery_order_details/:order_id", HarddeliveryOrderDetailsByOrderId)
	config.Router.OPTIONS("/harddeliverys/", HarddeliverysOptions)
	log.Println("Initialize Harddeliverys Runter ... ")
}

func HarddeliverysOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func HarddeliveryNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var harddeliverys Harddeliverys
	json.Unmarshal(body, &harddeliverys)

	var stocks stock.Stocks
	row := stocks.GetByID(harddeliverys.ID)

	// verify exist
	if len(row) == 0 {
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.ID_NOT_EXIST, StatusText: config.StatusText.ID_NOT_EXIST}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("ERROR: Insert record to tbl_harddeliverys failed, the record not exist! id:", harddeliverys.ID)
		return
	}

	sqlStatement := fmt.Sprintf("INSERT INTO tbl_harddeliverys (guid,order_id,date,id,price,project_id,soft_model,version,algorithm,period,user,delivery_type,tracking_number,description)VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')",
		config.Guid(), harddeliverys.Order_id, harddeliverys.Date, harddeliverys.ID, harddeliverys.Price, harddeliverys.Project_id, harddeliverys.Soft_model, harddeliverys.Version, harddeliverys.Algorithm, harddeliverys.Period, harddeliverys.User, harddeliverys.Delivery_type, harddeliverys.Tracking_number, harddeliverys.Description)
	log.Println(sqlStatement)

	//add new record to tbl_stocks
	stocks = stock.Stocks{Status: config.Business.DELIVERY, Date: harddeliverys.Date, Description: harddeliverys.Description, Event: sqlStatement}
	returned := stocks.UpdateStatusById(harddeliverys.ID)

	var returnObj config.ReturnStruct
	json.Unmarshal(returned, &returnObj)

	if returnObj.Status == config.StatusCode.UPDATE_NOT_ALLOW {
		fmt.Fprintln(w, string(returned))
		return
	}

	id, _ := gosql.Insert(config.DB, sqlStatement)
	log.Println("returned rowid=", id)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ = json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HarddeliveryRejected(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var harddeliverys Harddeliverys
	json.Unmarshal(body, &harddeliverys)

	var stocks stock.Stocks
	row := stocks.GetByID(harddeliverys.ID)

	// verify exist
	if len(row) == 0 {
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.ID_NOT_EXIST, StatusText: config.StatusText.ID_NOT_EXIST}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("ERROR: Insert record to tbl_harddeliverys failed, the record not exist! id:", harddeliverys.ID)
		return
	}

	//在生成出库单时需要通过注释来区分出库和退货
	harddeliverys.Description = "退货：" + harddeliverys.Description
	sqlStatement := fmt.Sprintf("INSERT INTO tbl_harddeliverys (guid,order_id,date,id,price,project_id,soft_model,version,algorithm,period,user,delivery_type,tracking_number,description)VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')",
		config.Guid(), harddeliverys.Order_id, harddeliverys.Date, harddeliverys.ID, harddeliverys.Price, harddeliverys.Project_id, harddeliverys.Soft_model, harddeliverys.Version, harddeliverys.Algorithm, harddeliverys.Period, harddeliverys.User, harddeliverys.Delivery_type, harddeliverys.Tracking_number, harddeliverys.Description)
	log.Println(sqlStatement)

	//add new record to tbl_stocks
	stocks = stock.Stocks{Status: config.Business.REJECT, Date: harddeliverys.Date, Description: harddeliverys.Description, Event: sqlStatement}
	returned := stocks.UpdateStatusById(harddeliverys.ID)

	var returnObj config.ReturnStruct
	json.Unmarshal(returned, &returnObj)

	if returnObj.Status == config.StatusCode.UPDATE_NOT_ALLOW {
		fmt.Fprintln(w, string(returned))
		return
	}

	id, _ := gosql.Insert(config.DB, sqlStatement)
	log.Println("returned id=", id)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ = json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HarddeliverysGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var harddeliverys Harddeliverys
	json.Unmarshal(body, &harddeliverys)
	condition := harddeliverys.convStructToString("AND")

	var sqlStatement string

	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("SELECT * FROM view_harddeliverys WHERE %s", condition)
	} else {
		// no condtion for select
		sqlStatement = fmt.Sprintf("SELECT * FROM view_harddeliverys")
	}

	log.Println(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HarddeliverysGetByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT * FROM view_harddeliverys WHERE guid='%s'", ps.ByName("guid"))
	log.Printf(sqlStatement)
	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HarddeliverysUpdateByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var harddeliverys Harddeliverys
	json.Unmarshal(body, &harddeliverys)

	condition := harddeliverys.convStructToString(",")
	var sqlStatement string
	if condition != "" {
		// have condition for select
		//多条数据，判断id和date，如果同一天多条数据，可能存在问题
		sqlStatement = fmt.Sprintf("UPDATE tbl_harddeliverys SET %s WHERE guid = '%s'", condition, ps.ByName("guid"))

		log.Println(sqlStatement)

		rowsAffected, _ := gosql.Update(config.DB, sqlStatement)

		responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowsAffected:%d}", sqlStatement, rowsAffected))
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("returned rowsAffected=", rowsAffected)

	} else {
		log.Printf("ERROR: No updated column specified. %s\n", sqlStatement)

		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.COLUMN_NOT_SPECIFIED, StatusText: config.StatusText.COLUMN_NOT_SPECIFIED}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))
	}

	// add event table for update hardstore column unless platform/model
	sqlStatement2 := fmt.Sprintf("SELECT id FROM tbl_harddeliverys WHERE guid='%s'", ps.ByName("guid"))
	log.Println(sqlStatement2)
	rows, _ := gosql.Query(config.DB, sqlStatement2)

	events := event.Events{
		ID:          (*rows)[0]["id"],
		Time:        time.Now().Format("2006-01-02 15:04:05"),
		Business:    config.Business.DELIVERY,
		Action:      config.Action.UPDATE,
		Operator:    config.OPERATOR,
		Description: sqlStatement,
	}
	events.Insert()
}

func HarddeliveryOrderList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	//出库单要包含退货的记录
	sqlStatement := fmt.Sprintf("SELECT DISTINCT order_id FROM view_harddeliverys")
	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HarddeliveryOrderDetailsByOrderId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	//这里只返回出库的记录，但不能按照状态判断，因为状态一旦改变，在生成历史出库单的时候会有问题
	sqlStatement := fmt.Sprintf("SELECT *, COUNT(*) as count FROM view_harddeliverys WHERE order_id='%s' ", ps.ByName("order_id"))
	//%在go语言中属于特殊字符，无法直接在格式化字符串中原样输出
	sqlStatement = sqlStatement + "AND description NOT LIKE '退货：%' GROUP BY model"
	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	for i := range *rows {
		sqlStatement = fmt.Sprintf("SELECT id FROM view_harddeliverys WHERE order_id='%s' AND model='%s' ", ps.ByName("order_id"), ((*rows)[i])["model"])
		sqlStatement = sqlStatement + "AND description NOT LIKE '退货：%'"
		id_list, _ := gosql.Query(config.DB, sqlStatement)
		id_list_json, _ := json.Marshal(id_list)
		((*rows)[i])["idlist"] = string(id_list_json)
	}

	//这里返回退货的记录，但不能按照状态判断，因为状态一旦改变，在生成历史出库单的时候会有问题
	sqlStatement2 := fmt.Sprintf("SELECT *, COUNT(*) as count FROM view_harddeliverys WHERE order_id='%s' ", ps.ByName("order_id"))
	sqlStatement2 = sqlStatement2 + "AND description LIKE '退货：%' GROUP BY model"
	log.Printf(sqlStatement2)

	rows_r, _ := gosql.Query(config.DB, sqlStatement2)
	for i := range *rows_r {
		sqlStatement2 = fmt.Sprintf("SELECT id FROM view_harddeliverys WHERE order_id='%s' AND model='%s' ", ps.ByName("order_id"), ((*rows_r)[i])["model"])
		sqlStatement2 = sqlStatement2 + "AND description LIKE '退货：%'"
		id_list, _ := gosql.Query(config.DB, sqlStatement2)
		id_list_json, _ := json.Marshal(id_list)
		((*rows_r)[i])["idlist"] = string(id_list_json)
		(*rows_r)[i]["count"] = "-" + (*rows_r)[i]["count"]
	}

	deliveryResponseText, _ := json.Marshal(rows)
	dlength := len(string(deliveryResponseText))

	rejectResponseText, _ := json.Marshal(rows_r)
	rlength := len(string(rejectResponseText))

	var responseText string
	// 将出库和退货的信息合并在一起返回给客户端
	// lengh <= 2 means empty
	if dlength <= 2 && rlength <= 2 {
		responseText = string(deliveryResponseText)
	} else if rlength <= 2 {
		responseText = string(deliveryResponseText)
	} else if dlength <= 2 {
		responseText = string(rejectResponseText)
	} else {
		responseText = "[" + (string(deliveryResponseText))[1:dlength-1] + "," + (string(rejectResponseText))[1:rlength-1] + "]"
	}

	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: responseText, Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

// construct select condition or update columns
func (harddeliverys *Harddeliverys) convStructToString(separator string) string {

	var condition string

	// don't update ID when PUT request
	if separator != "," && harddeliverys.ID != "" {
		condition += fmt.Sprintf("id='%s'%s", harddeliverys.ID, separator)
	}
	if harddeliverys.Order_id != "" {
		condition += fmt.Sprintf("order_id='%s'%s", harddeliverys.Order_id, separator)
	}
	if harddeliverys.Date != "" {
		condition += fmt.Sprintf("date='%s'%s", harddeliverys.Date, separator)
	}
	if harddeliverys.Price != "" {
		condition += fmt.Sprintf("price='%s'%s", harddeliverys.Price, separator)
	}
	if harddeliverys.Project_id != "" {
		condition += fmt.Sprintf("project_id='%s'%s", harddeliverys.Project_id, separator)
	}
	if harddeliverys.Soft_model != "" {
		condition += fmt.Sprintf("soft_model='%s'%s", harddeliverys.Soft_model, separator)
	}
	if harddeliverys.Version != "" {
		condition += fmt.Sprintf("version='%s'%s", harddeliverys.Version, separator)
	}
	if harddeliverys.Algorithm != "" {
		condition += fmt.Sprintf("algorithm='%s'%s", harddeliverys.Algorithm, separator)
	}
	if harddeliverys.Period != "" {
		condition += fmt.Sprintf("period='%s'%s", harddeliverys.Period, separator)
	}
	if harddeliverys.User != "" {
		condition += fmt.Sprintf("user='%s'%s", harddeliverys.User, separator)
	}
	if harddeliverys.Delivery_type != "" {
		condition += fmt.Sprintf("delivery_type='%s'%s", harddeliverys.Delivery_type, separator)
	}
	if harddeliverys.Tracking_number != "" {
		condition += fmt.Sprintf("tracking_number='%s'%s", harddeliverys.Tracking_number, separator)
	}
	if harddeliverys.Description != "" {
		condition += fmt.Sprintf("description='%s'%s", harddeliverys.Description, separator)
	}

	return strings.Trim(condition, separator)
}
