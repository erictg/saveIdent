package redisHelper

import (
	"testing"
	"log"
	"time"
)

func TestInit(t *testing.T) {
	Init("localhost:6379")

	Push("YOOOOO1")
	Push("Yooooooo 2")

	time.Sleep(time.Second + 3)

	str := redisClient.RPop("inputQueue")
	log.Println(str.String())

	str = redisClient.RPop("inputQueue")
	log.Println(str.String())


}