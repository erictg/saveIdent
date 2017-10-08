package main

import (
	"saveIdent.com/server/sqlInterfaceService/queries"
	"log"
)

func main(){

	//test user shit
	userQueries, err := queries.NewUserQueries()
	if err != nil{
		log.Fatal(err.Error())
	}

	err = userQueries.Create("Eric Solender", "dfadsf", "asdfasdf", "asdfasd")
	if err != nil{
		log.Fatal(err.Error())
	}

	user, err := userQueries.GetUserWithoutRolesAndPermissionByName("Eric Solender")
	if err != nil{
		log.Fatal(err.Error())
	}
	log.Println(user)
}