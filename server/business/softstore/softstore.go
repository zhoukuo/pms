package softstore

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

type Softstores struct {
	Model        string
	Version      string
	Count        string
	Publish_Date string
	Store_Date   string
	Class        string
	Description  string
}

func init() {
	// init router.
	config.Router.POST("/softstores/", SoftstoreNew)
	config.Router.GET("/softstores/", SoftstoresGet)
	config.Router.GET("/softstores/:guid", SoftstoreGetByGUID)
	config.Router.POST("/softstores/:guid", SoftstoreUpdateByGUID)
	config.Router.GET("/softstores_version_list/:model", SoftReturnVersionByModel)
	config.Router.OPTIONS("/softstores/", SoftstoresOptions)
	log.Println("Initialize Softstores Runter ... ")
}

func SoftstoresOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func SoftstoreNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var softstores Softstores
	json.Unmarshal(body, &softstores)

	sqlStatement := fmt.Sprintf("INSERT INTO tbl_softstores(guid,model,version,count,publish_date,store_date,class,description) VALUES('%s','%s','%s','%s','%s','%s','%s','%s')",
		config.Guid(), softstores.Model, softstores.Version, softstores.Count, softstores.Publish_Date, softstores.Store_Date, softstores.Class, softstores.Description)
	log.Println(sqlStatement)

	id, _ := gosql.Insert(config.DB, sqlStatement)
	log.Println("returned rowid=", id)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func SoftstoresGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var softstores Softstores
	json.Unmarshal(body, &softstores)
	condition := softstores.convStructToString("AND")

	var sqlStatement string
	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_softstores WHERE %s", condition)
	} else {
		// no condtion for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_softstores")
	}

	log.Println(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func SoftstoreGetByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT * FROM tbl_softstores WHERE guid='%s'", ps.ByName("guid"))

	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func SoftstoreUpdateByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var softstores Softstores
	json.Unmarshal(body, &softstores)

	condition := softstores.convStructToString(",")
	var sqlStatement string
	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("UPDATE tbl_softstores SET %s WHERE guid='%s'", condition, ps.ByName("guid"))
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

		log.Printf("ERROR: No upcountd column specified. %s\n", sqlStatement)
	}

}

func SoftReturnVersionByModel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT version FROM tbl_softstores WHERE model = '%s'", ps.ByName("model"))
	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

// construct select condition or upcount columns
func (softstores *Softstores) convStructToString(separator string) string {

	var condition string

	// don't upcount Model when PUT request
	if separator != "," && softstores.Model != "" {
		condition += fmt.Sprintf("model='%s'%s", softstores.Model, separator)
	}
	if softstores.Version != "" {
		condition += fmt.Sprintf("version='%s'%s", softstores.Version, separator)
	}
	if softstores.Count != "" {
		condition += fmt.Sprintf("count='%s'%s", softstores.Count, separator)
	}
	if softstores.Publish_Date != "" {
		condition += fmt.Sprintf("publish_date='%s'%s", softstores.Publish_Date, separator)
	}
	if softstores.Store_Date != "" {
		condition += fmt.Sprintf("store_date='%s'%s", softstores.Store_Date, separator)
	}
	if softstores.Class != "" {
		condition += fmt.Sprintf("class='%s'%s", softstores.Class, separator)
	}
	if softstores.Description != "" {
		condition += fmt.Sprintf("description='%s'%s", softstores.Description, separator)
	}

	return strings.Trim(condition, separator)
}
