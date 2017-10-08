package rest

import (
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"saveIdent.com/server/sqlInterfaceService/rest/dto"
	"saveIdent.com/server/sqlInterfaceService/queries"
	"saveIdent.com/common/httpHelper"
	dto2 "saveIdent.com/server/deviceInputService/dto"
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
		httpHelper.BuildJsonResponse(dto.ErrorResponse{Reason:err.Error()}, &res, http.StatusBadRequest)
		return
	}
	httpHelper.BuildJsonResponse(dto2.NoResp{}, &res, 201)
}

func FindUserHandler(res http.ResponseWriter, req *http.Request){
	q := req.URL.Query().Get("name")
	if q == ""{
		httpHelper.BuildJsonResponse(dto.ErrorResponse{Reason:"no query"}, &res, http.StatusBadRequest)
		return
	}

	user, err := queries.User.GetUserWithoutRolesAndPermissionByName(q)
	if err != nil{
		log.Println(err.Error())
		httpHelper.BuildJsonResponse(dto.ErrorResponse{Reason:err.Error()}, &res, http.StatusBadRequest)
		return
	}
	httpHelper.BuildJsonResponse(user, &res, 200)
}

func AddDeviceHandler(res http.ResponseWriter, req *http.Request){
	queries.Device.Create()
	httpHelper.BuildJsonResponse(dto2.NoResp{}, &res, http.StatusNoContent)
}

func AssociateDeviceWithUserHandler(res http.ResponseWriter, req *http.Request){
	reqDto := dto.AssociateDTO{}
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

	err = queries.Device.AssociateWithUser(reqDto.DeviceId, reqDto.UserId)
	if err != nil{
		log.Println("something happened")
		log.Println(err.Error())
		httpHelper.BuildJsonResponse(dto.ErrorResponse{Reason:err.Error()}, &res, http.StatusBadRequest)
		return
	}

	httpHelper.BuildJsonResponse(dto2.NoResp{}, &res, http.StatusCreated)
}

