package service

import (
	"saveIdent.com/common/models"
	"github.com/go-redis/redis"
	"time"
	"fmt"
	"saveIdent.com/common/httpHelper"
	"log"
	"encoding/json"
)

var client *redis.Client

func Init(addr string){
	client = redis.NewClient(&redis.Options{
		Addr:addr,
		Password:"",
		DB:0,
	})
}

func AddOrUpdate(location models.Location){
	bytes, _ := httpHelper.EncodeJson(location)
	client.Set(fmt.Sprintf("%d", location.DeviceId), string(bytes), time.Hour * 24)
}

func RemoveDevice(id uint){
	client.Del(fmt.Sprintf("%d", id))
}

func Get(id uint) (models.Location, error){
	result := client.Get(fmt.Sprintf("%d", id))
	if res, err := result.Result(); err != nil{
		log.Println(err.Error())
		return models.Location{}, err
	}else{
		reqDto := models.Location{}

		err = json.Unmarshal([]byte(res), &reqDto)
		if err != nil{
			log.Println("failed to unmarshall json")
			log.Println(err.Error())
			return models.Location{}, err
		}

		return reqDto, nil
	}
}

