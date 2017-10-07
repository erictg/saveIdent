package elasticService

import (
	"github.com/Zaba505/golang-generic-pool/src"
	"saveIdent.com/common/logger"
	"net/http"
	"time"
	"saveIdent.com/server/deviceInputService/dto"
	"encoding/json"
	"fmt"
	"bytes"
	"errors"
	"os"
)

const (
	UPDATE = "/shit/updateDtos"
)

type ElasticSearchDB struct {
	dbIp string
	clientPool *genericPool.GenericPool
	errLogger *logger.JsonLogger
}

func New(ipAddr string, numOfClients int, errLogger *logger.JsonLogger) *ElasticSearchDB {
	clientPool := genericPool.New(numOfClients, http.Client{}, func(fishieNum int) map[int]interface{} { return map[int]interface{}{4: time.Second * 30} })

	if errLogger == nil {
		errLogger = logger.NewLogger(os.Stdout, "  ")
	}

	return &ElasticSearchDB{ipAddr, clientPool, errLogger}
}

func (db *ElasticSearchDB) ChangeIPTo(newIp string) {
	db.dbIp = newIp
}

type Location struct {
	Type string		`json:"type"`
	Cors []float32	`json:"coordinates"`
}

type GeoLocation struct {
	Loc Location	`json:"location"`
}

type UpdateESDocType struct {
	DeviceId int		`json:"device_id"`
	Status int			`json:"status"`
	Geo GeoLocation		`json:"geo"`
}

func (db *ElasticSearchDB) Add(update dto.UpdateRequestDTO) error {

	client, ok := db.clientPool.GetFishie().(http.Client)
	if !ok {
		if db.errLogger != nil {
			db.errLogger.Warn(1, "No free fishie")
		}
		fmt.Println("No free fishie")
		return errors.New("no free fishie")
	}

	defer db.clientPool.PutFishieBack(client)

	loc := Location{"point", []float32{update.Lon, update.Lat}}
	geoLoc := GeoLocation{loc}
	shitToAdd := UpdateESDocType{update.Device_id, update.Status, geoLoc}

	b, err := json.Marshal(shitToAdd)
	if err != nil {
		if db.errLogger != nil {
			db.errLogger.Error(err, "Failed marshalling json")
		} else {
			fmt.Println("Failed marshalling json")
			panic(err)
		}
	}

	resp, err := client.Post(db.dbIp + UPDATE, "application/json", bytes.NewBuffer(b))
	defer resp.Body.Close()

	if err != nil {
		if db.errLogger != nil {
			db.errLogger.Error(err, "Failed POST of info to db")
		}
		fmt.Println(err)
		return err
	}

	return nil

}