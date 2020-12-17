package rest

import (
	"awesomeProject/peresistence"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type eventServiceHandler struct {
	dbhandler peresistence.DatabaseHandler
}

func NewEventHandler(databasehandler peresistence.DatabaseHandler) *eventServiceHandler {
	return &eventServiceHandler {
		dbhandler: databasehandler,
	}
}

func (eh * eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{error: No search criteria found, 
                   you can either search by id via /id/4
                   to search by name via /name/coldplayconcert}`)
		return
	}
	searchkey, ok := vars["Search"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{error: No search keys found, you can either search by id via /id/4
                   to search by name via /name/coldplayconcert}`)
		return
	}
	var event peresistence.Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbhandler.FindEventByName(searchkey)
	case "id":
		id, err := hex.DecodeString(searchkey)
		if err == nil {
			event, err = eh.dbhandler.FindEvent(id)
		}
	}
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "{error occured %s}", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh * eventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllAvailableEvents()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error: %s}", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&events)
}

func (eh * eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
	event := peresistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{error : %s}", err)
		return
	}
	id, err := eh.dbhandler.AddEvent(event)
	if nil != err {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{id: %d, error: %s}", id, err)
		return
	}
}