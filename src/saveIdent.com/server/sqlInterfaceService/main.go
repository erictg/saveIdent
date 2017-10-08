package main

import (
	"saveIdent.com/server/sqlInterfaceService/queries"
	"github.com/gorilla/mux"
	"saveIdent.com/server/sqlInterfaceService/rest"
	"fmt"
	"net/http"
	"github.com/gorilla/handlers"
)

func main(){
	queries.Init()

	router := mux.NewRouter()

	router.HandleFunc("/rest/user/register", rest.RegisterHandler).Methods("POST")
	router.HandleFunc("/rest/user", rest.FindUserHandler).Methods("GET")
	router.HandleFunc("/rest/device/new", rest.AddDeviceHandler).Methods("POST")
	router.HandleFunc("/rest/device/associate", rest.AssociateDeviceWithUserHandler).Methods("PUT")

	port := 1991
	portStr := fmt.Sprintf("%d", port)

	http.ListenAndServe(":"+portStr, handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With"}),
		handlers.IgnoreOptions())(router))
}

