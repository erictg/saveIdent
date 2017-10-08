package connectionES

import (
	"saveIdent.com/common/elasticService/connection"
	"saveIdent.com/common/logger"
	"os"
	"bufio"
)

var ElasticSearch *elasticService.ElasticSearchDB

func Init(addr string){
	logger := logger.NewLogger(bufio.NewWriter(os.Stdout), "	")
	ElasticSearch = elasticService.New(addr, 2, logger)
}