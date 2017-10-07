package connectionES

import (
	"saveIdent.com/common/elasticService/connection"
	"saveIdent.com/common/logger"
	"os"
	"bufio"
)

var ElasticSearch *connection.ElasticSearchDB

func Init(addr string){
	logger := logger.NewLogger(bufio.NewWriter(os.Stdout), "	")
	ElasticSearch = connection.New(addr, 2, logger)
}