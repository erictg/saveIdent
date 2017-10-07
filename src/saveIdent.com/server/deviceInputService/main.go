package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"saveIdent.com/server/deviceInputService/handler"
	"fmt"
	"net/http"
)

func main(){

	//registry.Init()

	router := mux.NewRouter()

	router.HandleFunc("/rest/device/update", handler.HandlePositionUpdate).Methods("POST")
	router.HandleFunc("/rest/device/status", handler.HandleStatusChange).Methods("POST")

	port := 1990
	portStr := fmt.Sprintf("%d", port)

	http.ListenAndServe(":"+portStr, handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With"}),
		handlers.IgnoreOptions())(router))
}
