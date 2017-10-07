package logger

import (
	"encoding/json"
	"io"
	"strconv"
	"runtime/debug"
	"time"
	"strings"
)

/*
	This package contains a logger specifically meant for JSON
	This file contains the logging functionality
*/

const (
	LOGS_BUFFER = 20
)

type Info struct {
	Type string
	Time string
	Message string
	Stack string
}

type Error struct {
	Type string
	Time string
	Message string
	Error string
	Stack string
}

type Warn struct {
	Type string
	Time string
	Level int
	Message string
	Stack string
}

func (log Info) String() string { return "{ Type: INFO, Message: " + log.Message  + " }" }

func (log Warn) String() string { return "{ Type: WARN, Level: " + strconv.Itoa(log.Level) + ", Message: " + log.Message + "}" }

func (log Error) String() string { return "{ Type: ERROR, Message: " + log.Message + ", Error: " + log.Error + "}" }

type Logger interface {
	Info(message string)
	Warn(level int, message string)
	Error(err error, message string)
	Log()
}

type JsonLogger struct {
	encoder *json.Encoder
	shitToLog chan interface{}
}

func NewLogger(w io.Writer, indent string) *JsonLogger {

	shitToLogChan := make(chan interface{}, LOGS_BUFFER)

	encoder := json.NewEncoder(w)

	encoder.SetIndent("", indent)

	jsonLogger := &JsonLogger{encoder, shitToLogChan}

	go jsonLogger.Log()

	return jsonLogger

}

func (logger *JsonLogger) Log() {
	for thangToLog := range logger.shitToLog {
		logger.encoder.Encode(thangToLog)
	}
}

func (logger *JsonLogger) Info(message string) {

	logger.shitToLog <- &Info{"INFO", time.Now().String(), message, getStackTrace()}

}

func (logger *JsonLogger) Warn(level int, message string) {

	logger.shitToLog <- &Warn{"WARN", time.Now().String(), level, message, getStackTrace()}

}

func (logger *JsonLogger) Error(err error, message string) {

	logger.shitToLog <- &Error{"ERROR", time.Now().String(), message, err.Error(), getStackTrace()}

}


// This gets the stack trace of the goroutine that logs sumthing
// It ignores the calls to debug.Stack() and logger.getStackTrace() in its output
func getStackTrace() string {

	stackString := string(debug.Stack())

	return strings.Join(strings.Split(stackString,"\n")[5:], "\n")
}