package handler

import (
	"net/http"
	"saveIdent.com/server/deviceInputService/dto"
	"saveIdent.com/common/httpHelper"
	"log"
	"io/ioutil"
	"encoding/json"
)

//needs to cache into ES
//TODO make ES cache service
//dump into proc queue
//check if the location is in an affected area
func HandlePositionUpdate(res http.ResponseWriter, req *http.Request){
	reqDto := dto.UpdateRequestDTO{}
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
	if err != nil{
		httpHelper.BuildJsonResponse(dto.NoResp{}, &res, 500)
		return
	}

	log.Println(reqDto)

	//todo dump into elastic search

	//todo dump into the queue

	//todo wait for check

	httpHelper.BuildJsonResponse(dto.UpdateResponseDTO{IsAffected:false}, &res, http.StatusOK)
}

func HandleStatusChange(res http.ResponseWriter, req *http.Request){
	reqDto := dto.StatusChangedRequestDTO{}
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

	//todo do something with this change

	log.Println(reqDto)

	httpHelper.BuildJsonResponse(dto.NoResp{}, &res, http.StatusNoContent)
}