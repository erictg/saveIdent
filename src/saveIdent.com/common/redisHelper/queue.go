package redisHelper

import (


	"github.com/go-redis/redis"
	"log"
	"saveIdent.com/server/deviceInputService/dto"
)

var pushQueue chan dto.UpdateRequestDTO
var redisClient *redis.Client

func Init(addr string){
	redisClient = redis.NewClient(&redis.Options{
		Addr:addr,
		Password:"",
		DB:0,
	})
	pushQueue = make(chan  dto.UpdateRequestDTO, 100)

	go PushRoutine(pushQueue)
}

func Push(string dto.UpdateRequestDTO){
	pushQueue <- string
}

func PushRoutine(queue chan dto.UpdateRequestDTO){
	log.Println("inside of the routine")
	for{

		select {
		case str := <-queue:
			log.Println("pushing to redis!")
			log.Println(str)
			redisClient.RPush("inputLatQueue", str.Lat)
			redisClient.RPush("inputLonQueue", str.Lon)
			redisClient.RPush("inputIdQueue", str.Device_id)
		}


	}
}


