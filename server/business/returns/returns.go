package returns

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

type Returns struct {
	ID          string
	Date        string
	Description string
}

func init() {
	// init router.
	config.Router.POST("/returns/", ReturnNew)
	config.Router.GET("/returns/", ReturnsGet)
	config.Router.GET("/returns/:guid", ReturnGetByGUID)
	config.Router.POST("/returns/:guid", ReturnUpdateByGUID)
	config.Router.OPTIONS("/returns/", ReturnsOptions)
	log.Println("Initialize Returns Runter ... ")
}

func ReturnsOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func ReturnNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var returns Returns
	json.Unmarshal(body, &returns)

	var stocks stock.Stocks
	row := stocks.GetByID(returns.ID)

	// verify exist
	if len(row) == 0 {
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.ID_NOT_EXIST, StatusText: config.StatusText.ID_NOT_EXIST}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("ERROR: Insert record to tbl_returns failed, the record not exist! id:", returns.ID)
		return
	}

	sqlStatement := fmt.Sprintf("INSERT INTO tbl_returns(guid, id, date, description) VALUES('%s','%s','%s', '%s')",
		config.Guid(), returns.ID, returns.Date, returns.Description)
	log.Println(sqlStatement)

	//update tbl_stocks status
	stocks = stock.Stocks{Status: config.Business.RETURN, Date: returns.Date, Description: returns.Description, Event: sqlStatement}
	returned := stocks.UpdateStatusById(returns.ID)

	var returnObj config.ReturnStruct
	json.Unmarshal(returned, &returnObj)

	if returnObj.Status == config.StatusCode.UPDATE_NOT_ALLOW {
		fmt.Fprintln(w, string(returned))
		return
	}

	id, _ := gosql.Insert(config.DB, sqlStatement)
	log.Println("returned", id)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ = json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func ReturnsGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var returns Returns
	json.Unmarshal(body, &returns)
	condition := returns.convStructToString("AND")

	var sqlStatement string
	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("SELECT * FROM view_returns WHERE %s", condition)
	} else {
		// no condtion for select
		sqlStatement = fmt.Sprintf("SELECT * FROM view_returns")
	}
	log.Println(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func ReturnGetByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT * FROM view_returns WHERE guid='%s'", ps.ByName("guid"))

	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func ReturnUpdateByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var returns Returns
	json.Unmarshal(body, &returns)

	condition := returns.convStructToString(",")
	var sqlStatement string
	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("UPDATE tbl_returns SET %s WHERE guid='%s'", condition, ps.ByName("guid"))
		log.Println(sqlStatement)

		rowsAffected, _ := gosql.Update(config.DB, sqlStatement)

		responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowsAffected:%d}", sqlStatement, rowsAffected))
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("returned", rowsAffected)

	} else {
		log.Printf("ERROR: No updated column specified. %s\n", sqlStatement)

		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.COLUMN_NOT_SPECIFIED, StatusText: config.StatusText.COLUMN_NOT_SPECIFIED}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))
	}

	sqlStatement2 := fmt.Sprintf("SELECT id FROM view_returns WHERE guid='%s'", ps.ByName("guid"))
	log.Println(sqlStatement2)
	rows, _ := gosql.Query(config.DB, sqlStatement2)

	events := event.Events{
		ID:          (*rows)[0]["id"],
		Time:        time.Now().Format("2006-01-02 15:04:05"),
		Business:    config.Business.RETURN,
		Action:      config.Action.UPDATE,
		Operator:    config.OPERATOR,
		Description: sqlStatement,
	}
	events.Insert()
}

// construct select condition or update columns
func (returns *Returns) convStructToString(separator string) string {

	var condition string

	// don't update ID when PUT request
	if separator != "," && returns.ID != "" {
		condition += fmt.Sprintf("id='%s'%s", returns.ID, separator)
	}
	if returns.Date != "" {
		condition += fmt.Sprintf("date='%s'%s", returns.Date, separator)
	}
	if returns.Description != "" {
		condition += fmt.Sprintf("description='%s'%s", returns.Description, separator)
	}

	return strings.Trim(condition, separator)
}
