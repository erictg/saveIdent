package main

import (
	"saveIdent.com/server/activeMonitorService/service"
	"github.com/gorilla/mux"
	"saveIdent.com/server/activeMonitorService/rest"
	"fmt"
	"net/http"
	"github.com/gorilla/handlers"
)

func main(){
	service.Init("localhost:6731")
	router := mux.NewRouter()

	router.HandleFunc("/rest/monitor", rest.HandleAddUpdate).Methods("POST")
	router.HandleFunc("/rest/monitor", rest.HandleRemoveDevice).Methods("DELETE")
	router.HandleFunc("/rest/monitor", rest.HandleGetDevice).Methods("GET")
	port := 1992
	portStr := fmt.Sprintf("%d", port)

	http.ListenAndServe(":"+portStr, handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With"}),
		handlers.IgnoreOptions())(router))
}