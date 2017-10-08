package rest

import (
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"saveIdent.com/common/models"
	"saveIdent.com/server/activeMonitorService/service"
	"saveIdent.com/common/httpHelper"
	"saveIdent.com/server/deviceInputService/dto"
	"strconv"
	dto2 "saveIdent.com/server/sqlInterfaceService/rest/dto"
)

func HandleAddUpdate(res http.ResponseWriter, req *http.Request){
	reqDto := models.Location{}
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

	service.AddOrUpdate(reqDto)

	httpHelper.BuildJsonResponse(dto.NoResp{}, &res, http.StatusNoContent)
}

func HandleRemoveDevice(res http.ResponseWriter, req *http.Request){
	id := req.URL.Query().Get("id")
	i, err := strconv.ParseUint(id, 10, 32)
	if err != nil{
		httpHelper.BuildJsonResponse(dto2.ErrorResponse{Reason:"int is fucked"}, &res, 500)
	}
	service.RemoveDevice(uint(i))
	httpHelper.BuildJsonResponse(dto.NoResp{}, &res, http.StatusNoContent)

}

func HandleGetDevice(res http.ResponseWriter, req *http.Request){
	id := req.URL.Query().Get("id")
	i, err := strconv.ParseUint(id, 10, 32)
	if err != nil{
		httpHelper.BuildJsonResponse(dto2.ErrorResponse{Reason:"int is fucked"}, &res, 500)
	}
	loc, err := service.Get(uint(i))
	if err != nil{
		log.Fatal(err.Error())
	}
	httpHelper.BuildJsonResponse(loc, &res, http.StatusNoContent)
}