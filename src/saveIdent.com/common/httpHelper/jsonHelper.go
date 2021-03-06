package httpHelper

import (
	"net/http"
	"encoding/json"
//	"io/ioutil"
	"log"
	"fmt"
)

func EncodeJson(encode interface{}) ([]byte, error){
	return json.Marshal(encode)
}

func DecodeJsonFromString(string string, copyTo interface{}) error{
	return DecodeJson([]byte(string), copyTo)
}

func DecodeJson(bytes []byte, copyTo interface{}) error{
	return json.Unmarshal(bytes, copyTo)
}
////decodes json from request object into a holder interface/struct
//func DecodeJsonFromRequest(r *http.Request) (type{}, error){
//	var copyTo interface{}
//	bytes, err := ioutil.ReadAll(r.Body)
//	if err != nil{
//		log.Println("something fucked up")
//		log.Println(err.Error())
//		return nil, err
//	}
//
//
//	err = DecodeJson(bytes, &copyTo)
//	if err != nil{
//		log.Println("failed to unmarshall json")
//		log.Println(err.Error())
//		return nil, err
//	}
//
//	return copyTo, nil
//}



func BuildJsonResponse(encode interface{}, response *http.ResponseWriter, httpStatus int) error{

	body, err := EncodeJson(encode)
	if err != nil{
		log.Println("failed to encode")
		(*response).WriteHeader(http.StatusInternalServerError)
		return err
	}
	(*response).WriteHeader(httpStatus)
	(*response).Header().Add("Accept", "application/json")
	fmt.Fprintf(*response, "%s", body)

	return nil
}