package hardware

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

type Hardwares struct {
	Platform     string
	Supplier     string
	Type         string
	Unit         string
	Price        string
	Class        string
	Publish_date string
	Description  string
}

func init() {
	// init router.
	config.Router.POST("/hardwares/", HardwareNew)
	config.Router.GET("/hardwares/", HardwaresGet)
	config.Router.GET("/hardwares/:guid", HardwareGetByGUID)
	config.Router.POST("/hardwares/:guid", HardwareUpdateByGUID)
	config.Router.OPTIONS("/hardwares/", HardwaresOptions)
	log.Println("Initialize Hardwares Runter ... ")
}

func HardwaresOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func HardwareNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var hardwares Hardwares
	json.Unmarshal(body, &hardwares)

	//check if it exist in database
	sqlStatement := fmt.Sprintf("SELECT * FROM tbl_hardwares WHERE platform='%s'", hardwares.Platform)
	rows, _ := gosql.Query(config.DB, sqlStatement)
	if len(*rows) > 0 {
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.PLATFORM_EXIST, StatusText: config.StatusText.PLATFORM_EXIST}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("Insert record to tbl_hardwares failed, the record exist! platform:", hardwares.Platform)
		return
	}

	//insert data to tbl_loans
	sqlStatement = fmt.Sprintf("INSERT INTO tbl_hardwares(guid, platform, supplier, type, unit, price, class, publish_date, description) VALUES('%s','%s','%s','%s','%s','%s','%s','%s','%s')",
		config.Guid(), hardwares.Platform, hardwares.Supplier, hardwares.Type, hardwares.Unit, hardwares.Price, hardwares.Class, hardwares.Publish_date, hardwares.Description)
	log.Println(sqlStatement)
	//execute the sql
	id, _ := gosql.Insert(config.DB, sqlStatement)
	//print the result of executed
	log.Println("returned id=", id)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HardwaresGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var hardwares Hardwares
	json.Unmarshal(body, &hardwares)
	condition := hardwares.convStructToString("AND")

	var sqlStatement string

	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_hardwares WHERE %s", condition)
	} else {
		// no condtion for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_hardwares")
	}

	log.Println(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HardwareGetByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT * FROM tbl_hardwares WHERE guid='%s'", ps.ByName("guid"))

	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func HardwareUpdateByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var hardwares Hardwares
	json.Unmarshal(body, &hardwares)

	condition := hardwares.convStructToString(",")
	var sqlStatement string
	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("UPDATE tbl_hardwares SET %s WHERE guid='%s'", condition, ps.ByName("guid"))
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

}

// construct select condition or update columns
func (hardwares *Hardwares) convStructToString(separator string) string {

	var condition string

	// don't update ID when PUT request
	if separator != "," && hardwares.Platform != "" {
		condition += fmt.Sprintf("platform='%s'%s", hardwares.Platform, separator)
	}
	if hardwares.Supplier != "" {
		condition += fmt.Sprintf("supplier='%s'%s", hardwares.Supplier, separator)
	}
	if hardwares.Type != "" {
		condition += fmt.Sprintf("type='%s'%s", hardwares.Type, separator)
	}
	if hardwares.Unit != "" {
		condition += fmt.Sprintf("unit='%s'%s", hardwares.Unit, separator)
	}
	if hardwares.Price != "" {
		condition += fmt.Sprintf("price='%s'%s", hardwares.Price, separator)
	}
	if hardwares.Class != "" {
		condition += fmt.Sprintf("class='%s'%s", hardwares.Class, separator)
	}
	if hardwares.Publish_date != "" {
		condition += fmt.Sprintf("publish_date='%s'%s", hardwares.Publish_date, separator)
	}
	if hardwares.Description != "" {
		condition += fmt.Sprintf("description='%s'%s", hardwares.Description, separator)
	}

	return strings.Trim(condition, separator)
}
