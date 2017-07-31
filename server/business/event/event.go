package event

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

type Events struct {
	ID          string
	Time        string
	Business    string
	Action      string
	Operator    string
	Description string
}

func init() {
	// init router.
	config.Router.POST("/events/", EventNew)
	config.Router.GET("/events/", EventsGet)
	config.Router.GET("/events/:guid", EventsGetByGUID)
	config.Router.OPTIONS("/events/", EventsOptions)
	log.Println("Initialize Events Runter ... ")
}

func EventsOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)
}

func EventNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var events Events
	json.Unmarshal(body, &events)

	res := events.Insert()
	fmt.Fprintln(w, res)
}

func EventsGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	var events Events
	json.Unmarshal(body, &events)
	condition := events.convStructToString("AND")

	var sqlStatement string

	if condition != "" {
		// have condition for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_events WHERE %s", condition)
	} else {
		// no condtion for select
		sqlStatement = fmt.Sprintf("SELECT * FROM tbl_events")
	}

	log.Println(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

func EventsGetByGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	config.SetAccessControlAllowOrigin(w, r)

	sqlStatement := fmt.Sprintf("SELECT * FROM tbl_events WHERE guid='%s'", ps.ByName("guid"))

	log.Printf(sqlStatement)

	rows, _ := gosql.Query(config.DB, sqlStatement)

	responseText, _ := json.Marshal(rows)
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	fmt.Fprintln(w, string(returned))
}

// construct select condition or update columns
func (events *Events) convStructToString(separator string) string {

	var condition string

	// don't update ID when PUT request
	if separator != "," && events.ID != "" {
		condition += fmt.Sprintf("id='%s'%s", events.ID, separator)
	}
	if events.Time != "" {
		condition += fmt.Sprintf("time='%s'%s", events.Time, separator)
	}
	if events.Business != "" {
		condition += fmt.Sprintf("business='%s'%s", events.Business, separator)
	}
	if events.Action != "" {
		condition += fmt.Sprintf("action='%s'%s", events.Action, separator)
	}
	if events.Operator != "" {
		condition += fmt.Sprintf("operator='%s'%s", events.Operator, separator)
	}
	if events.Description != "" {
		condition += fmt.Sprintf("description='%s'%s", events.Description, separator)
	}

	return strings.Trim(condition, separator)
}

func (events *Events) Insert() (writer string) {

	sqlStatement := fmt.Sprintf("INSERT INTO tbl_events(guid, id, time, business, action, operator, description) VALUES('%s','%s','%s','%s','%s','%s',\"%s\")",
		config.Guid(), events.ID, events.Time, events.Business, events.Action, events.Operator, events.Description)
	log.Println(sqlStatement)

	id, _ := gosql.Insert(config.DB, sqlStatement)

	responseText, _ := json.Marshal(fmt.Sprintf("{sqlStatement:%s, rowid:%d}", sqlStatement, id))
	var ReturnObject = config.ReturnStruct{ReadyStatus: 4, ResponseText: string(responseText), Status: 200, StatusText: "OK"}
	returned, _ := json.Marshal(ReturnObject)
	writer = string(returned)

	log.Println("returned rowid=", id)
	return
}
