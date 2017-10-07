package registry

import (
	"github.com/BurntSushi/toml"
	"log"
)

var configRead = false

var Data Globals

func Init(){
	InitRegistryFile("config.toml")
}

func InitRegistryFile(configFile string){
	if !configRead{
		configRead = true
		if _, err := toml.Decode(configFile, &Data); err != nil{
			log.Println(err.Error())
			log.Fatal("failed to load config file")
		}
	}
}


type Globals struct{
	Host string
	Port int
	ClientId string
}