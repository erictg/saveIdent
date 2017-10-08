package rest

import (
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"saveIdent.com/server/sqlInterfaceService/rest/dto"
	"saveIdent.com/server/sqlInterfaceService/queries"
	"saveIdent.com/common/httpHelper"
)

//handle user shit
func RegisterHandler(res http.ResponseWriter, req *http.Request){
	reqDto := dto.RegisterRequestDTO{}
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil{
		log.Println("something fucked up")
		log.Println(err.Error())
		return
	}


	err = json.Unmarshal(bytes, &reqDto)
	if err != nil{
		log.Println("failed to unmarshall json")
		log.Println(err.Error())
	}

	err = queries.User.Create(reqDto.Name, reqDto.Password, reqDto.Email, reqDto.Address)
	if err != nil{
		log.Println(err)
		log.Println("failed to create user")
		httpHelper.BuildJsonResponse(dto.ErrorResponse{Reason:err.Error()}, &res, 400)
	}
}
