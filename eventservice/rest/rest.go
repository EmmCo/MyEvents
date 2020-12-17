package rest

import (
	"awesomeProject/peresistence"
	"github.com/gorilla/mux"
	"net/http"
)

func ServeAPI(endpoint string, dbHandler peresistence.DatabaseHandler) error {
	handler := NewEventHandler(dbHandler)
	rtr := mux.NewRouter()
	eventsrouter := rtr.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.allEventHandler)
	return http.ListenAndServe(endpoint, rtr)
}
