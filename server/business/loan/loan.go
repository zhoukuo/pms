package loan

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

type Loans struct {
	ID              string
	Date            string
	Project_id      string
	User            string
	Delivery_type   string
	Tracking_number string
	Description     string
}

func init() {
	// init router.
	config.Router.POST("/loans/", LoanNew)
	config.Router.GET("/loans/", LoansGet)
	config.Router.GET("/loans/:guid", LoanGetByGUID)
	config.Router.POST("/loans/:guid", LoanUpdateByGUID)
	config.Router.OPTIONS("/loans/", LoansOptions)
	log.Println("Initialize Loans Runter ... ")
}

func LoansOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func LoanNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var loans Loans
	json.Unmarshal(body, &loans)

	var stocks stock.Stocks
	row := stocks.GetByID(loans.ID)

	// verify exist
	if len(row) == 0 {
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.ID_NOT_EXIST, StatusText: config.StatusText.ID_NOT_EXIST}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("ERROR: Insert record to tbl_loans failed, the record not exist! id:", loans.ID)
		return
	}

	//insert data to tbl_loans
	sqlStatement := fmt.Sprintf("INSERT INTO tbl_loans(guid, id, date, project_id, user, delivery_type, tracking_number, description) VALUES('%s','%s','%s','%s','%s','%s','%s', '%s')",
		config.Guid(), loans.ID, loans.Date, loans.Project_id, loans.User, loans.Delivery_type, loans.Tracking_number, loans.Description)
	log.Println(sqlStatement)

	//update tbl_stocks status
	stocks = stock.Stocks{Status: config.Business.LOAN, Date: loans.Date, Description: loans.Description, Event: sqlStatement}
	returned := stocks.UpdateStatusById(loans.ID)

	var returnObj config.ReturnStruct
	json.Unmarshal(returned, &returnObj)

	if returnObj.Status == config.StatusCode.UPDATE_NOT_ALLOW {
		fmt.Fprintln(w, string(returned))
		return
	}

	//execute the sql
	id, _ := gosql.Insert(config.DB, sqlStatement)
	//print the result of executed
	log.Println("returned id=", id)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ = json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func LoansGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var loans Loans
	json.Unmarshal(body, &loans)
	condition := loans.convStructToString("AND")

	var sqlStatement string

	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("SELECT * FROM view_loans WHERE %s", condition)
	} else {
		// no condtion for select
		sqlStatement = fmt.Sprintf("SELECT * FROM view_loans")
	}

	log.Println(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func LoanGetByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT * FROM view_loans WHERE guid='%s'", ps.ByName("guid"))

	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func LoanUpdateByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var loans Loans
	json.Unmarshal(body, &loans)

	condition := loans.convStructToString(",")
	var sqlStatement string
	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("UPDATE tbl_loans SET %s WHERE guid='%s'", condition, ps.ByName("guid"))
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

	sqlStatement2 := fmt.Sprintf("SELECT id FROM view_loans WHERE guid='%s'", ps.ByName("guid"))
	log.Println(sqlStatement2)
	rows, _ := gosql.Query(config.DB, sqlStatement2)

	events := event.Events{
		ID:          (*rows)[0]["id"],
		Time:        time.Now().Format("2006-01-02 15:04:05"),
		Business:    config.Business.LOAN,
		Action:      config.Action.UPDATE,
		Operator:    config.OPERATOR,
		Description: sqlStatement,
	}
	events.Insert()
}

// construct select condition or update columns
func (loans *Loans) convStructToString(separator string) string {

	var condition string

	// don't update ID when PUT request
	if separator != "," && loans.ID != "" {
		condition += fmt.Sprintf("id='%s'%s", loans.ID, separator)
	}
	if loans.Date != "" {
		condition += fmt.Sprintf("date='%s'%s", loans.Date, separator)
	}
	if loans.Project_id != "" {
		condition += fmt.Sprintf("project_id='%s'%s", loans.Project_id, separator)
	}
	if loans.User != "" {
		condition += fmt.Sprintf("user='%s'%s", loans.User, separator)
	}
	if loans.Delivery_type != "" {
		condition += fmt.Sprintf("delivery_type='%s'%s", loans.Delivery_type, separator)
	}
	if loans.Tracking_number != "" {
		condition += fmt.Sprintf("tracking_number='%s'%s", loans.Tracking_number, separator)
	}
	if loans.Description != "" {
		condition += fmt.Sprintf("description='%s'%s", loans.Description, separator)
	}

	return strings.Trim(condition, separator)
}
