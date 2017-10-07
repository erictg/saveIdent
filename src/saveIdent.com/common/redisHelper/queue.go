package redisHelper

import (


	"github.com/go-redis/redis"
	"log"
)

var pushQueue chan string
var redisClient *redis.Client

func Init(addr string){
	redisClient = redis.NewClient(&redis.Options{
		Addr:addr,
		Password:"",
		DB:0,
	})
	pushQueue = make(chan  string, 100)

	go PushRoutine(pushQueue)
}

func Push(string string){
	pushQueue <- string
}

func PushRoutine(queue chan string){
	log.Println("inside of the routine")
	for{

		select {
		case str := <-queue:
			log.Println("pushing to redis!")
			log.Println(str)
			redisClient.RPush("inputQueue", str)
		}


	}
}


type queue struct{
	q []string
}

func(q *queue) Init(){
	q.q = make([]string, 0)
}

func (q *queue) Push(str string){
	q.q = append(q.q, str)
}
//
//func Pop(q chan []string) (string, error){
//	if .IsEmpty(){
//		return "", errors.New("queue is empty")
//	}
//	toReturn := q.q[0]
//	q.q = q.q[1:]
//	return toReturn, nil
//}

func (q *queue) IsEmpty() bool{
	return len(q.q) == 0
}