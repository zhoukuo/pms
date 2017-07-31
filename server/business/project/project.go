package project

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

type Projects struct {
	ID          string
	Name        string
	Manager     string
	Dept        string
	Customer    string
	Description string
}

func init() {
	// init router.
	config.Router.POST("/projects/", ProjectNew)
	config.Router.GET("/projects/", ProjectsGet)
	config.Router.GET("/projects/:guid", ProjectsGetByGUID)
	config.Router.POST("/projects/:guid", ProjectsUpdateByGUID)
	config.Router.OPTIONS("/projects/", ProjectsOptions)
	log.Println("Initialize Projects Runter ... ")
}

func ProjectsOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func ProjectNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
	body, _ := ioutil.ReadAll(r.Body)
	var projects Projects
	// all request must be json format
	json.Unmarshal(body, &projects)

	//check if it exist in database
	sqlStatement := fmt.Sprintf("SELECT * FROM tbl_projects WHERE id='%s'", projects.ID)
	rows, _ := gosql.Query(config.DB, sqlStatement)
	if len(*rows) > 0 {
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.PROJECT_EXIST, StatusText: config.StatusText.PROJECT_EXIST}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("Insert record to tbl_projects failed, the record exist! id:", projects.ID)
		return
	}

	sqlStatement = fmt.Sprintf("INSERT INTO tbl_projects(guid,id,name,manager,dept,customer,description) VALUES('%s','%s','%s','%s','%s','%s','%s')",
		config.Guid(), projects.ID, projects.Name, projects.Manager, projects.Dept, projects.Customer, projects.Description)
	log.Println(sqlStatement)

	id, _ := gosql.Insert(config.DB, sqlStatement)
	log.Println("returned rowid=", id)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func ProjectsGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
	body, _ := ioutil.ReadAll(r.Body)
	var projects Projects
	json.Unmarshal(body, &projects)
	condition := projects.convStructToString("AND")

	var sqlStatement string

	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_projects WHERE %s", condition)
	} else {
		// no condtion for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_projects")
	}

	log.Println(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func ProjectsGetByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT * FROM tbl_projects WHERE guid='%s'", ps.ByName("guid"))
	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func ProjectsUpdateByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
	body, _ := ioutil.ReadAll(r.Body)
	var projects Projects
	json.Unmarshal(body, &projects)

	condition := projects.convStructToString(",")
	var sqlStatement string
	if condition != "" {
		sqlStatement = fmt.Sprintf("UPDATE tbl_projects SET %s WHERE guid='%s'", condition, ps.ByName("guid"))
		log.Println(sqlStatement)

		rowsAffected, _ := gosql.Update(config.DB, sqlStatement)

		responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowsAffected:%d}", sqlStatement, rowsAffected))
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Println("returned rowsAffected=", rowsAffected)

	} else {
		var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: "null", Status: config.StatusCode.COLUMN_NOT_SPECIFIED, StatusText: config.StatusText.COLUMN_NOT_SPECIFIED}
		returned, _ := json.Marshal(ReturnObject)
		fmt.Fprintln(w, string(returned))

		log.Printf("ERROR: No updated column specified. %s\n", sqlStatement)
	}
}

// construct select condition or update columns
func (projects *Projects) convStructToString(separator string) string {
	var condition string
	// don't update ID when PUT request
	if separator != "," && projects.ID != "" {
		condition += fmt.Sprintf("id='%s'%s", projects.ID, separator)
	}
	if projects.Name != "" {
		condition += fmt.Sprintf("name='%s'%s", projects.Name, separator)
	}
	if projects.Manager != "" {
		condition += fmt.Sprintf("manager='%s'%s", projects.Manager, separator)
	}
	if projects.Dept != "" {
		condition += fmt.Sprintf("dept='%s'%s", projects.Dept, separator)
	}
	if projects.Customer != "" {
		condition += fmt.Sprintf("customer='%s'%s", projects.Customer, separator)
	}
	if projects.Description != "" {
		condition += fmt.Sprintf("description='%s'%s", projects.Description, separator)
	}

	return strings.Trim(condition, separator)
}
