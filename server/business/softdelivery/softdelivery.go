package softdelivery

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/zhoukuo/gosql"
	"io/ioutil"
	"log"
	"net/http"
	"pms/config"
	"strings"
)

type Softdeliverys struct {
	Date          string
	Model         string
	Count         string
	Version       string
	Project_ID    string
	Delivery_Type string
	Tracking_Number   string
	Description   string
}

func init() {
	//init router
	config.Router.POST("/softdeliverys/", SoftdeliveryNew)
	config.Router.GET("/softdeliverys/", SoftdeliveryGet)
	config.Router.GET("/softdeliverys/:guid", GetSoftdeliveryByGUID)
	config.Router.POST("/softdeliverys/:guid", SoftdeliveryUpdateByGUID)
	config.Router.OPTIONS("/softdeliverys/", SoftdeliverysOptions)
	log.Println("Initialize softdelivery Runter ... ")
}

func SoftdeliverysOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func SoftdeliveryNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var softdeliverys Softdeliverys
	// all request must be json format
	json.Unmarshal(body, &softdeliverys)

	sqlStatement := fmt.Sprintf("INSERT INTO tbl_softdeliverys(guid,date,model,count,version,project_id,delivery_type,tracking_number,description) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s')",
		config.Guid(), softdeliverys.Date, softdeliverys.Model, softdeliverys.Count, softdeliverys.Version, softdeliverys.Project_ID, softdeliverys.Delivery_Type, softdeliverys.Tracking_Number, softdeliverys.Description)
	log.Println(sqlStatement)

	id, _ := gosql.Insert(config.DB, sqlStatement)
	log.Println("returned rowid=", id)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func SoftdeliveryGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var softdeliverys Softdeliverys
	json.Unmarshal(body, &softdeliverys)
	condition := softdeliverys.convStructToString("AND")

	var sqlStatement string

	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_softdeliverys WHERE %s", condition)
	} else {
		// no condtion for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_softdeliverys")
	}

	log.Println(sqlStatement)
	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func GetSoftdeliveryByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT * FROM tbl_softdeliverys WHERE guid='%s'", ps.ByName("guid"))

	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))

}

func SoftdeliveryUpdateByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var softdeliverys Softdeliverys
	json.Unmarshal(body, &softdeliverys)

	condition := softdeliverys.convStructToString(",")

	var sqlStatement string
	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("UPDATE tbl_softdeliverys SET %s WHERE guid='%s'", condition, ps.ByName("guid"))
		log.Println(sqlStatement)

		rowsAffected, _ := gosql.Update(config.DB, sqlStatement)
		log.Println("returned rowsAffected=", rowsAffected)

		responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowsAffected:%d}", sqlStatement, rowsAffected))
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

	} else {
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.COLUMN_NOT_SPECIFIED, StatusText: config.StatusText.COLUMN_NOT_SPECIFIED}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Printf("ERROR: No updated column specified. %s\n", sqlStatement)
	}

}
func (softdeliverys *Softdeliverys) convStructToString(separator string) string {
	// construct select condition or update columns

	var condition string

	// don't update ID when PUT request
	if separator != "," && softdeliverys.Project_ID != "" {
		condition += fmt.Sprintf("project_id='%s'%s", softdeliverys.Project_ID, separator)
	}
	if softdeliverys.Model != "" {
		condition += fmt.Sprintf("model='%s'%s", softdeliverys.Model, separator)
	}
	if softdeliverys.Count != "" {
		condition += fmt.Sprintf("count='%s'%s", softdeliverys.Count, separator)
	}
	if softdeliverys.Version != "" {
		condition += fmt.Sprintf("version='%s'%s", softdeliverys.Version, separator)
	}
	if softdeliverys.Date != "" {
		condition += fmt.Sprintf("date='%s'%s", softdeliverys.Date, separator)
	}
	if softdeliverys.Delivery_Type != "" {
		condition += fmt.Sprintf("delivery_type='%s'%s", softdeliverys.Delivery_Type, separator)
	}
	if softdeliverys.Tracking_Number != "" {
		condition += fmt.Sprintf("tracking_number='%s'%s", softdeliverys.Tracking_Number, separator)
	}
	if softdeliverys.Description != "" {
		condition += fmt.Sprintf("description='%s'%s", softdeliverys.Description, separator)
	}

	return strings.Trim(condition, separator)
}
