package logger

import (
	"testing"
	"os"
	"time"
	"errors"
	"fmt"
)

var logger *JsonLogger
var parser *Parser

func TestMain(m *testing.M) {

	testFile, err := os.OpenFile("testFile.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		panic(err)
	}

	defer testFile.Close()

	logger = NewLogger(testFile, "  ")
	parser = NewParser(testFile)
	defer parser.Close()

	os.Exit(m.Run())

}


/* Here are all the logger.Parser tests */

func TestParser_ExtractLog(t *testing.T) {

	receiveLogChan := make(chan ParsedLog)
	receiveRequest := ExtractRequest{1, receiveLogChan}

	log := parser.ExtractLog(receiveRequest)

	if log.Err != nil {
		t.Error(log.Err)
	}

	switch t := log.Log.(type) {
	case Info:
		PrettyPrint(t)
	case Warn:
		PrettyPrint(t)
	case Error:
		PrettyPrint(t)
	default:
		fmt.Println(t, log)
	}

}

func TestParser_ExtractLogs(t *testing.T) {

	logs := parser.ExtractLogs()

	for _, log := range logs {

		if log.Err != nil {
			t.Error(log.Err)
		}

		switch t := log.Log.(type) {
		case Info:
			PrettyPrint(t)
		case Warn:
			PrettyPrint(t)
		case Error:
			PrettyPrint(t)
		default:
			fmt.Println(t, log)
		}
	}
}


/* Here are all the logger.JsonLogger tests */

func TestJsonLogger_Info(t *testing.T) {

	logger.Info("INFO test")

	time.Sleep(time.Second * 2)

}

func TestJsonLogger_Warn(t *testing.T) {

	logger.Warn(3, "WARN test")

	time.Sleep(time.Second * 2)

}

func TestJsonLogger_Error(t *testing.T) {

	logger.Error(errors.New("New error for testing"), "ERROR test")

	time.Sleep(time.Second * 2)

}