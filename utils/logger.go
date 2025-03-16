package utils

import (
	"log"
	"os"
)

const loggerFlags = log.Ldate | log.Ltime | log.Lshortfile

type CustomLogger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

var instance *CustomLogger

func GetCustomLogger() *CustomLogger {

	if instance == nil {
		instance = &CustomLogger{
			InfoLogger:  log.New(os.Stdout, "[INFO]\t", loggerFlags),
			ErrorLogger: log.New(os.Stderr, "[ERROR]\t", loggerFlags),
		}
	}

	return instance
}
