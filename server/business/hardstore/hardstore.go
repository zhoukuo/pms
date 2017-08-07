package hardstore

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

type Hardstores struct {
	ID          string
	Order_id    string
	Date        string
	Price       string
	SAP         string
	Receiving   string
	Description string
	Model       string
	Platform    string
	Status      string
}

func init() {
	// init router.
	config.Router.POST("/hardstores/", HardstoreNew)
	config.Router.GET("/hardstores/", HardstoresGet)
	config.Router.GET("/hardstores/:guid", HardstoreGetByGUID)
	config.Router.POST("/hardstores/:guid", HardstoreUpdateByGUID)
	config.Router.GET("/hardstore_orderlist/", HardstoreOrderList)
	config.Router.GET("/hardstore_order_details/:order_id", HardstoreOrderDetailsByOrderId)
	config.Router.OPTIONS("/hardstores/", HardstoresOptions)
	log.Println("Initialize Hardstores Runter ... ")
}

func HardstoresOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func HardstoreNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var hardstores Hardstores
	json.Unmarshal(body, &hardstores)

	// verify status of this record if allow insert new record
	var stocks stock.Stocks
	if stocks.ValidateEnabled(hardstores.ID, config.Business.STORE) == false {
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.ID_ALREADY_EXIST, StatusText: config.StatusText.ID_ALREADY_EXIST}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("Insert record to tbl_hardstores failed, the record exist! id:", hardstores.ID)
		return
	}

	sqlStatement := fmt.Sprintf("INSERT INTO tbl_hardstores(guid,order_id,date,id,price,sap,receiving, description) VALUES('%s','%s','%s','%s','%s','%s','%s','%s')",
		config.Guid(), hardstores.Order_id, hardstores.Date, hardstores.ID, hardstores.Price, hardstores.SAP, hardstores.Receiving, hardstores.Description)
	log.Println(sqlStatement)

	id, _ := gosql.Insert(config.DB, sqlStatement)
	log.Println("returned rowid=", id)

	//入库事件的描述信息和其他业务不同，因为入库的主要信息都保存在库存表中，因此这里的事件信息为空，在库存表中添加记录时添加描述
	stocks = stock.Stocks{
		ID:          hardstores.ID,
		Platform:    hardstores.Platform,
		Model:       hardstores.Model,
		Status:      config.Business.STORE,
		Date:        hardstores.Date,
		Description: hardstores.Description,
		Event:       "",
	}
	stocks.Insert()

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HardstoresGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var hardstores Hardstores
	json.Unmarshal(body, &hardstores)
	condition := hardstores.convStructToString("AND")

	var sqlStatement string

	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("SELECT * FROM view_hardstores WHERE %s", condition)
	} else {
		// no condtion for select
		sqlStatement = fmt.Sprintf("SELECT * FROM view_hardstores")
	}

	log.Println(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HardstoreGetByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT * FROM view_hardstores WHERE guid='%s'", ps.ByName("guid"))
	log.Printf(sqlStatement)
	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HardstoreUpdateByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var hardstores Hardstores
	json.Unmarshal(body, &hardstores)
	condition := hardstores.convStructToString(",")

	sqlStatement := fmt.Sprintf("SELECT id FROM tbl_hardstores WHERE guid='%s'", ps.ByName("guid"))
	log.Println(sqlStatement)
	rows, _ := gosql.Query(config.DB, sqlStatement)

	// update platform/model in tbl_stocks and set status=作废
	if hardstores.Status == config.Business.INVALID {

		// verify status of this record if allow invalid
		var stocks stock.Stocks
		if stocks.ValidateEnabled((*rows)[0]["id"], config.Business.INVALID) == false {
			var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.INVALID_NOT_ALLOW, StatusText: config.StatusText.INVALID_NOT_ALLOW}
			returned, _ := json.Marshal(ReturnObject)
			fmt.Fprintln(w, string(returned))

			log.Println("ERROR: Update status to invalid failed, because status is not '入库' guid:", ps.ByName("guid"))
			return
		}

		// 作废的时间不是入库的时间，而是当前操作的时间
		stocks = stock.Stocks{Status: hardstores.Status, Date: time.Now().Format("2006-01-02"), Description: hardstores.Description, Event: hardstores.Description}
		returned := stocks.UpdateStatusById((*rows)[0]["id"])
		fmt.Fprintln(w, string(returned))

		//作废的记录直接删除
		sqlStatement = fmt.Sprintf("DELETE FROM tbl_hardstores WHERE guid='%s'", ps.ByName("guid"))
		log.Println(sqlStatement)
		rowsAffected, _ := gosql.Delete(config.DB, sqlStatement)
		log.Println("returned rowsAffected=", rowsAffected)

		//
		sqlStatement = fmt.Sprintf("DELETE FROM tbl_stocks WHERE id='%s'", (*rows)[0]["id"])
		log.Println(sqlStatement)
		rowsAffected, _ = gosql.Delete(config.DB, sqlStatement)
		log.Println("returned rowsAffected=", rowsAffected)
		//这里主要是操作库存表，因此不返回入库表的状态
		return
	}

	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("UPDATE tbl_hardstores SET %s WHERE guid='%s'", condition, ps.ByName("guid"))
		log.Println(sqlStatement)

		rowsAffected, _ := gosql.Update(config.DB, sqlStatement)
		log.Println("returned rowsAffected=", rowsAffected)

		responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowsAffected:%d}", sqlStatement, rowsAffected))
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

	} else {
		log.Printf("ERROR: No updated column specified. %s\n", sqlStatement)

		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.COLUMN_NOT_SPECIFIED, StatusText: config.StatusText.COLUMN_NOT_SPECIFIED}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))
	}

	//update platform/model in tbl_stocks
	stocks := stock.Stocks{Platform: hardstores.Platform, Model: hardstores.Model, Date: hardstores.Date}
	stocks.UpdateById((*rows)[0]["id"])

	events := event.Events{
		ID:          (*rows)[0]["id"],
		Time:        time.Now().Format("2006-01-02 15:04:05"),
		Business:    config.Business.STORE,
		Action:      config.Action.UPDATE,
		Operator:    config.OPERATOR,
		Description: sqlStatement}
	events.Insert()
}

func HardstoreOrderList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	//排除作废的入库记录
	sqlStatement := fmt.Sprintf("SELECT DISTINCT order_id FROM view_hardstores WHERE status<>'作废'")
	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HardstoreOrderDetailsByOrderId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	//排除作废的入库记录
	sqlStatement := fmt.Sprintf("SELECT *, COUNT(*) as count FROM view_hardstores WHERE status<>'作废' AND order_id='%s' GROUP BY platform", ps.ByName("order_id"))
	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

// construct select condition or update columns
func (hardstores *Hardstores) convStructToString(separator string) string {

	var condition string

	// don't update ID when PUT request
	if separator != "," && hardstores.ID != "" {
		condition += fmt.Sprintf("id='%s'%s", hardstores.ID, separator)
	}
	if hardstores.Order_id != "" {
		condition += fmt.Sprintf("order_id='%s'%s", hardstores.Order_id, separator)
	}
	if hardstores.Date != "" {
		condition += fmt.Sprintf("date='%s'%s", hardstores.Date, separator)
	}
	if hardstores.Price != "" {
		condition += fmt.Sprintf("price='%s'%s", hardstores.Price, separator)
	}
	if hardstores.SAP != "" {
		condition += fmt.Sprintf("sap='%s'%s", hardstores.SAP, separator)
	}
	if hardstores.Receiving != "" {
		condition += fmt.Sprintf("receiving='%s'%s", hardstores.Receiving, separator)
	}
	if hardstores.Description != "" {
		condition += fmt.Sprintf("description='%s'%s", hardstores.Description, separator)
	}

	return strings.Trim(condition, separator)
}
