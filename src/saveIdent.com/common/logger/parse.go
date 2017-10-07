package logger

import (
	"encoding/json"
	"io"
	"reflect"
	"os"
)

var INFO_STRUCT_TYPE = structType(Info{})
var WARN_STRUCT_TYPE = structType(Warn{})
var ERROR_STRUCT_TYPE = structType(Error{})

var INFO_TYPE = reflect.TypeOf(Info{})
var WARN_TYPE = reflect.TypeOf(Warn{})
var ERROR_TYPE = reflect.TypeOf(Error{})


/* This file dedicated to parsing the logs written by a JsonLogger */

type Parse interface {
	ExtractLog(receive chan ParsedLog) ParsedLog
	ExtractLogs() []ParsedLog
}


type Parser struct {
	decoders []*json.Decoder
	extractNew chan ExtractRequest
	quit chan bool
}

type ParsedLog struct {
	Log interface{}
	Err error
}

type ExtractRequest struct {
	ReaderNum int
	LogChan chan ParsedLog
}

// Given an io.Reader, which is pointing to where your logs are written to
// A new Parser is created around the reader
func NewParser(readers ...io.Reader) *Parser {

	extractNewLog, quit := make(chan ExtractRequest), make(chan bool)

	var decoders []*json.Decoder

	for _, r := range readers {
		decoders = append(decoders, json.NewDecoder(r))
	}

	parser := &Parser{decoders, extractNewLog, quit}

	go parser.extract()

	return parser

}

func (parser *Parser) extract() {
	for {
		select {

		case <-parser.quit:
			break

		case requestChan := <-parser.extractNew:

			var logVal map[string]interface{}

			err := parser.decoders[requestChan.ReaderNum - 1].Decode(&logVal)

			requestChan.LogChan <- ParsedLog{convert(logVal), err}

		}
	}
}

// This decodes exactly one log from the provided reader
// What ever happens to be the first one
func (parser *Parser) ExtractLog(receive ExtractRequest) ParsedLog {

	go func() { parser.extractNew <- receive }()

	return <-receive.LogChan

}

// This decodes all logs from the provided reader
func (parser *Parser) ExtractLogs() []ParsedLog {

	var logs []ParsedLog
	var log ParsedLog

	receiveLogChan := make(chan ParsedLog)
	receiveRequest := ExtractRequest{1, receiveLogChan}

	for parser.decoders[receiveRequest.ReaderNum - 1].More() {

		log = parser.ExtractLog(receiveRequest)

		if log.Err != nil {
			return logs
		}

		logs = append(logs, log)

	}

	return logs

}

func (parser *Parser) Close() {
	parser.quit <- true
}


// This is simply a helper function for printing logs out to Stdout
// in a pretty printed JSON format i.e. it has indentation
func PrettyPrint(logs ...interface{}) {

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	for _, log := range logs {
		if err := encoder.Encode(&log); err != nil {
			panic(err)
		}
	}

}


// This is a helper function for converting the decoded map into
// a log struct
func convert(logInterface map[string]interface{}) interface{} {

	var structFields []reflect.StructField
	var fieldType reflect.Type
	var warnLevelType int

	for field, val := range logInterface {

		fieldType = reflect.TypeOf(val)

		if fieldType.Kind() == reflect.Float64 {
			fieldType = reflect.TypeOf(warnLevelType)
		}

		structFields = append(structFields, reflect.StructField{Name: field, Type: fieldType})
	}

	logStruct := reflect.New(reflect.StructOf(structFields)).Elem()

	var infoFields, warnFields, errorFields int

	for i := 0; i < logStruct.NumField(); i++ {
		f := logStruct.Type().Field(i)

		switch {
		case containsField(INFO_STRUCT_TYPE, f):
			infoFields++
		case containsField(WARN_STRUCT_TYPE, f):
			warnFields++
		case containsField(ERROR_STRUCT_TYPE, f):
			errorFields++
		}
	}

	var newLogStruct reflect.Value
	var newLogStructFields []reflect.StructField
	var specificLogStruct interface{}

	switch {

	case infoFields | errorFields == 5:
		newLogStructFields = reArrangStructFields(ERROR_STRUCT_TYPE, logStruct)

		newLogStruct = reflect.New(reflect.StructOf(newLogStructFields)).Elem()

		setStructVals(logInterface, newLogStruct)

		specificLogStruct = newLogStruct.Convert(ERROR_TYPE).Interface()

	case infoFields | warnFields == 5:
		newLogStructFields = reArrangStructFields(WARN_STRUCT_TYPE, logStruct)

		newLogStruct = reflect.New(reflect.StructOf(newLogStructFields)).Elem()

		setStructVals(logInterface, newLogStruct)

		specificLogStruct = newLogStruct.Convert(WARN_TYPE).Interface()

	default:
		newLogStructFields = reArrangStructFields(INFO_STRUCT_TYPE, logStruct)

		newLogStruct = reflect.New(reflect.StructOf(newLogStructFields)).Elem()

		setStructVals(logInterface, newLogStruct)

		specificLogStruct = newLogStruct.Convert(INFO_TYPE).Interface()
	}

	return specificLogStruct

}


// Helper function for checking if a struct contains a field or not
func containsField(structType reflect.Type, field reflect.StructField) bool {
	_, ok := structType.FieldByName(field.Name)
	return ok
}


func reArrangStructFields(structType reflect.Type, oldStruct reflect.Value) []reflect.StructField {
	var newLogStructFields []reflect.StructField

	for j := 0; j < structType.NumField(); j++ {
		structField, _ := oldStruct.Type().FieldByName( structType.Field(j).Name )
		newLogStructFields = append(newLogStructFields, structField)
	}

	return newLogStructFields
}


func setStructVals(structFields map[string]interface{}, logStruct reflect.Value) {
	var warnLevelType int

	for field, val := range structFields {

		valVal := reflect.ValueOf(val)

		if valVal.Kind() == reflect.Float64 {
			valVal = valVal.Convert(reflect.TypeOf(warnLevelType))
		}

		logStruct.FieldByName(field).Set( valVal )
	}
}

// This is merely a helper funtion for declaring const values
func structType(logType interface{}) reflect.Type {

	var structFields []reflect.StructField

	logStructType := reflect.TypeOf(logType)

	for i := 0; i < logStructType.NumField(); i++ {
		structFields = append(structFields, logStructType.Field(i))
	}

	logStruct := reflect.New(reflect.StructOf(structFields)).Elem()

	return logStruct.Type()

}