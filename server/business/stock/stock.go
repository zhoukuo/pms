package stock

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/zhoukuo/gosql"
	"io/ioutil"
	"log"
	"net/http"
	"pms/business/event"
	"pms/config"
	"strings"
	"time"
)

type Stocks struct {
	ID          string
	Platform    string
	Model       string
	Status      string
	Date        string
	Description string
	Event       string // just for add event description
}

func init() {
	// init router.
	config.Router.POST("/stocks/", StockNew)
	config.Router.GET("/stocks/", StocksGet)
	config.Router.GET("/stocks/:id", StocksGetByID)
	config.Router.PUT("/stocks/:id", StocksUpdateByID)
	config.Router.GET("/stocks_useble/:status", StocksGetUsebleListByStatus)
	config.Router.GET("/stocks_check/:date", StocksCheckByDate)
	config.Router.OPTIONS("/stocks/", StocksOptions)
	log.Println("Initialize Stocks Runter ... ")
}

func StocksOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func StockNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
	body, _ := ioutil.ReadAll(r.Body)
	var stocks Stocks
	// all request must be json format
	json.Unmarshal(body, &stocks)

	res := stocks.Insert()
	fmt.Fprintln(w, res)
}

func StocksGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
	body, _ := ioutil.ReadAll(r.Body)
	var stocks Stocks
	json.Unmarshal(body, &stocks)
	condition := stocks.convStructToString("AND")

	var sqlStatement string

	if condition != "" {
		// have condition for select, 排除作废的入库记录
		sqlStatement = fmt.Sprintf("SELECT * FROM view_stocks WHERE status<>'作废' AND %s", condition)
	} else {
		// no condtion for select,排除作废的入库记录
		sqlStatement = fmt.Sprintf("SELECT * FROM view_stocks WHERE status<>'作废'")
	}

	log.Println(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func StocksGetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	var stocks Stocks
	row := stocks.GetByID(ps.ByName("id"))

	responseText, _ := json.Marshal(row)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func StocksUpdateByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
	body, _ := ioutil.ReadAll(r.Body)
	var stocks Stocks
	json.Unmarshal(body, &stocks)

	res := stocks.UpdateById(ps.ByName("id"))
	fmt.Fprintln(w, res)
}

func StocksGetUsebleListByStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	var stocks Stocks
	rows := stocks.GetUsebleListByStatus(ps.ByName("status"))

	var ReturnObject config.ReturnStruct
	if rows == nil {
		ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.STATUS_UNKNOWN, StatusText: config.StatusText.STATUS_UNKNOWN}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))
		return
	}

	responseText, _ := json.Marshal(rows)
	ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func StocksCheckByDate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	//排除作废/报废的记录
	sqlStatement := fmt.Sprintf("SELECT platform, SUM(case when status<>'出库' AND date(lastchanged_date)<date('%s','start of month') then 1 else 0 end) as count_lastmonth, SUM(case when date(store_date)>=date('%s', 'start of month') AND date(store_date)<=date('%s') then 1 else 0 end) as count_currentmonth_store, SUM(case when status='出库' AND date(delivery_date)>=date('%s', 'start of month') AND delivery_date<=date('%s') then 1 else 0 end) as count_currentmonth_delivery FROM view_stockscheck WHERE status<>'作废' AND status<>'报废' group by platform;", ps.ByName("date"), ps.ByName("date"), ps.ByName("date"), ps.ByName("date"), ps.ByName("date"))
	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	//rows中未包含入库单列表和出库单列表，需要单独增加
	for i := range *rows {
		//入库单列表
		sqlStatement = fmt.Sprintf("SELECT DISTINCT store_order_id as order_id FROM view_stockscheck WHERE date(store_date)>=date('%s', 'start of month') AND date(store_date)<=date('%s') AND status<>'作废' AND platform='%s'", ps.ByName("date"), ps.ByName("date"), (*rows)[i]["platform"])
		store_orderlist, _ := gosql.Query(config.DB, sqlStatement)
		store_orderlist_json, _ := json.Marshal(store_orderlist)
		((*rows)[i])["store_order_id"] = string(store_orderlist_json)

		//出库单列表
		sqlStatement = fmt.Sprintf("SELECT DISTINCT delivery_order_id as order_id FROM view_stockscheck WHERE date(delivery_date)>=date('%s', 'start of month') AND delivery_date<=date('%s') AND platform='%s'", ps.ByName("date"), ps.ByName("date"), (*rows)[i]["platform"])
		delivery_orderlist, _ := gosql.Query(config.DB, sqlStatement)
		delivery_orderlist_json, _ := json.Marshal(delivery_orderlist)
		((*rows)[i])["delivery_order_id"] = string(delivery_orderlist_json)
	}

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

// construct select condition or update columns
func (stocks *Stocks) convStructToString(separator string) string {
	var condition string
	// don't update ID when PUT request
	if separator != "," && stocks.ID != "" {
		condition += fmt.Sprintf("id='%s'%s", stocks.ID, separator)
	}
	if stocks.Platform != "" {
		condition += fmt.Sprintf("platform='%s'%s", stocks.Platform, separator)
	}
	if stocks.Model != "" {
		condition += fmt.Sprintf("model='%s'%s", stocks.Model, separator)
	}
	if stocks.Status != "" {
		condition += fmt.Sprintf("status='%s'%s", stocks.Status, separator)
	}
	if stocks.Description != "" {
		condition += fmt.Sprintf("description=\"%s\"%s", stocks.Description, separator)
	}

	return strings.Trim(condition, separator)
}

// This method for HardstoreNew() only!!!!
func (stocks *Stocks) Insert() (writer string) {
	sqlStatement := fmt.Sprintf("INSERT INTO tbl_stocks(id, platform, model, status, date, description) VALUES('%s','%s','%s','%s','%s', '%s')",
		stocks.ID, stocks.Platform, stocks.Model, stocks.Status, stocks.Date, stocks.Description)

	log.Println(sqlStatement)

	id, _ := gosql.Insert(config.DB, sqlStatement)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	writer = string(returned)

	log.Println("returned id=", id)

	// all insert must be add to event table
	description := fmt.Sprintf("平台:%s,  型号:%s,  日期:%s,  备注:%s", stocks.Platform, stocks.Model, stocks.Date, stocks.Description)
	events := event.Events{
		ID:          stocks.ID,
		Time:        time.Now().Format("2006-01-02 15:04:05"),
		Business:    stocks.Status,
		Action:      config.Action.CREATE,
		Operator:    config.OPERATOR,
		Description: description,
	}
	events.Insert()

	return
}

// This method for HardstoreUpdateByID() only!!!!
func (stocks *Stocks) UpdateById(id string) (writer string) {
	condition := stocks.convStructToString(",")
	var sqlStatement string
	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("UPDATE tbl_stocks SET %s WHERE id='%s'", condition, id)
		log.Println(sqlStatement)

		rowsAffected, _ := gosql.Update(config.DB, sqlStatement)
		log.Println("returned rowsAffected=", rowsAffected)

		responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowsAffected:%d}", sqlStatement, rowsAffected))
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
		returned, _ := json.Marshal(ReturnObject)
		writer = string(returned)

	} else {
		log.Printf("ERROR: No updated column specified. %s\n", sqlStatement)

		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.COLUMN_NOT_SPECIFIED, StatusText: config.StatusText.COLUMN_NOT_SPECIFIED}
		returned, _ := json.Marshal(ReturnObject)
		writer = string(returned)

		return
	}

	events := event.Events{
		ID:          id,
		Time:        time.Now().Format("2006-01-02 15:04:05"),
		Business:    config.Business.STORE,
		Action:      config.Action.UPDATE,
		Operator:    config.OPERATOR,
		Description: sqlStatement,
	}
	events.Insert()

	return
}

func (stocks *Stocks) UpdateStatusById(id string) (writer []byte) {
	// verify status of this record if allow update for this action
	if stocks.ValidateEnabled(id, stocks.Status) == false {
		log.Printf("ERROR: ValidateEnabled() return false! id:%s, status:%s\n", id, stocks.Status)

		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.UPDATE_NOT_ALLOW, StatusText: config.StatusText.UPDATE_NOT_ALLOW}
		writer, _ := json.Marshal(ReturnObject)
		return writer
	}

	sqlStatement := fmt.Sprintf("UPDATE tbl_stocks SET status='%s', date='%s', description='%s' WHERE id='%s'", stocks.Status, stocks.Date, stocks.Description, id)
	log.Println(sqlStatement)

	rowsAffected, _ := gosql.Update(config.DB, sqlStatement)
	log.Println("returned rowsAffected=", rowsAffected)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowsAffected:%d}", sqlStatement, rowsAffected))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	writer, _ = json.Marshal(ReturnObject)

	// all update status operation must be add to event table
	index := strings.LastIndex(stocks.Event, "VALUES")
	if index == -1 {
		index = 0
	}
	description := stocks.Event[index:]
	events := event.Events{
		ID:          id,
		Time:        time.Now().Format("2006-01-02 15:04:05"),
		Business:    stocks.Status,
		Action:      config.Action.CREATE,
		Operator:    config.OPERATOR,
		Description: description,
	}
	events.Insert()

	return writer
}

func (stocks *Stocks) GetByID(id string) map[string]string {

	sqlStatement := fmt.Sprintf("SELECT * FROM view_stocks WHERE id='%s'", id)
	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)
	if len(*rows) == 0 {
		log.Println("GetByID() no record returned!")
		// 如果记录不存在，返回一个空的map
		return make(map[string]string)
	}

	return (*rows)[0]
}

func (stocks *Stocks) ValidateEnabled(id string, status string) bool {
	row := stocks.GetByID(id)
	enable := false

	switch status {

	case config.Business.STORE:
		// 如果记录不存在，len(row) == 0
		enable = len(row) == 0

	case config.Business.INVALID:
		enable = row["status"] == config.Business.STORE

	case config.Business.DELIVERY:
		enable = row["status"] == config.Business.STORE || row["status"] == config.Business.RETURN || row["status"] == config.Business.REJECT || row["status"] == config.Business.FIXED

	case config.Business.REJECT:
		enable = row["status"] == config.Business.DELIVERY

	case config.Business.LOAN:
		enable = row["status"] == config.Business.STORE || row["status"] == config.Business.RETURN || row["status"] == config.Business.REJECT || row["status"] == config.Business.FIXED

	case config.Business.RETURN:
		enable = row["status"] == config.Business.LOAN

	case config.Business.REPAIR:
		enable = row["status"] == config.Business.STORE || row["status"] == config.Business.REJECT || row["status"] == config.Business.RETURN || row["status"] == config.Business.FIXED

	case config.Business.FIXED:
		enable = row["status"] == config.Business.REPAIR

	case config.Business.DISUSE:
		enable = row["status"] == config.Business.STORE || row["status"] == config.Business.REJECT || row["status"] == config.Business.RETURN || row["status"] == config.Business.FIXED
	default:
		enable = false
	}

	return enable
}

func (stocks *Stocks) GetUsebleListByStatus(status string) *[]map[string]string {
	var sqlStatement string

	switch status {

	case "STORE":
		// nothing to do
		sqlStatement = ""

	case "INVALID":
		sqlStatement = fmt.Sprintf("SELECT id, guid, status FROM tbl_stocks WHERE status='%s'", config.Business.STORE)

	case "REJECT":
		sqlStatement = fmt.Sprintf("SELECT id, guid, status FROM tbl_stocks WHERE status='%s'", config.Business.DELIVERY)

	case "RETURN":
		sqlStatement = fmt.Sprintf("SELECT id, guid, status FROM tbl_stocks WHERE status='%s'", config.Business.LOAN)

	case "FIXED":
		sqlStatement = fmt.Sprintf("SELECT id, guid, status FROM tbl_stocks WHERE status='%s'", config.Business.REPAIR)

	case "DELIVERY":
		sqlStatement = fmt.Sprintf("SELECT id, guid, status FROM tbl_stocks WHERE status='%s' OR status='%s' OR status='%s' OR status='%s'", config.Business.STORE, config.Business.REJECT, config.Business.RETURN, config.Business.FIXED)

	case "LOAN":
		sqlStatement = fmt.Sprintf("SELECT id, guid, status FROM tbl_stocks WHERE status='%s' OR status='%s' OR status='%s' OR status='%s'", config.Business.STORE, config.Business.REJECT, config.Business.RETURN, config.Business.FIXED)

	case "REPAIR":
		sqlStatement = fmt.Sprintf("SELECT id, guid, status FROM tbl_stocks WHERE status='%s' OR status='%s' OR status='%s' OR status='%s'", config.Business.STORE, config.Business.REJECT, config.Business.RETURN, config.Business.FIXED)

	case "DISUSE":
		sqlStatement = fmt.Sprintf("SELECT id, guid, status FROM tbl_stocks WHERE status='%s' OR status='%s' OR status='%s' OR status='%s'", config.Business.STORE, config.Business.REJECT, config.Business.RETURN, config.Business.FIXED)

	default:
		//  nothing to do
		sqlStatement = ""
	}

	if sqlStatement == "" {
		log.Println("ERROR: GetUsebleListByStatus() status is unknown!")
		// 状态未知，返回nil
		return nil
	}

	log.Println(sqlStatement)
	rows, _ := gosql.Query(config.DB, sqlStatement)
	return rows
}
