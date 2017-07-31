package software

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

type Softwares struct {
	Name        string
	Model       string
	Dept        string
	Description string
}

func init() {
	// init router.
	config.Router.POST("/softwares/", SoftwaresNew)
	config.Router.GET("/softwares/", SoftwaresGet)
	config.Router.GET("/softwares/:guid", SoftwaresGetByGUID)
	config.Router.POST("/softwares/:guid", SoftwaresUpdateByGUID)
	config.Router.OPTIONS("/softwares/", SoftwaresOptions)
	log.Println("Initialize Softwares Runter ... ")
}

func SoftwaresOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func SoftwaresNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var softwares Softwares
	// all request must be json format
	json.Unmarshal(body, &softwares)

	//check if it exist in database
	sqlStatement := fmt.Sprintf("SELECT * FROM tbl_softwares WHERE model='%s'", softwares.Model)
	rows, _ := gosql.Query(config.DB, sqlStatement)
	if len(*rows) > 0 {
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.MODEL_EXIST, StatusText: config.StatusText.MODEL_EXIST}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("Insert record to tbl_softwares failed, the record exist! model:", softwares.Model)
		return
	}

	sqlStatement = fmt.Sprintf("INSERT INTO tbl_softwares(guid, name, model, dept, description) VALUES('%s','%s','%s','%s','%s')",
		config.Guid(), softwares.Name, softwares.Model, softwares.Dept, softwares.Description)
	log.Println(sqlStatement)

	id, _ := gosql.Insert(config.DB, sqlStatement)
	log.Println("returned id=", id)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func SoftwaresGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var softwares Softwares
	json.Unmarshal(body, &softwares)
	condition := softwares.convStructToString("AND")

	var sqlStatement string

	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_softwares WHERE %s", condition)
	} else {
		// no condtion for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_softwares")
	}

	log.Println(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))

}

func SoftwaresGetByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT * FROM tbl_softwares WHERE guid='%s'", ps.ByName("guid"))
	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func SoftwaresUpdateByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var softwares Softwares
	json.Unmarshal(body, &softwares)
	condition := softwares.convStructToString(",")
	var sqlStatement string
	if condition != "" {
		sqlStatement = fmt.Sprintf("UPDATE tbl_softwares SET %s WHERE guid='%s'", condition, ps.ByName("guid"))
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
func (softwares *Softwares) convStructToString(separater string) string {
	var condition string

	if separater != "," && softwares.Model != "" {
		condition += fmt.Sprintf("model = '%s'%s", softwares.Model, separater)
	}

	if softwares.Name != "" {
		condition += fmt.Sprintf("name = '%s'%s", softwares.Name, separater)
	}

	if softwares.Dept != "" {
		condition += fmt.Sprintf("dept = '%s'%s", softwares.Dept, separater)
	}
	if softwares.Description != "" {
		condition += fmt.Sprintf("description = '%s'%s", softwares.Description, separater)
	}

	return strings.Trim(condition, separater)
}
