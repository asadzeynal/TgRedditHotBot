package util

import (
	"log"
	"os"
)

type Logger interface {
	Info(s string, v ...any)
	Warn(s string, v ...any)
	Error(s error, v ...any)
}

type CustomLog struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
}

func NewCustomLog() *CustomLog {
	infoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger := log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)

	logger := CustomLog{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		warnLogger:  warnLogger,
	}

	return &logger
}

func (logger *CustomLog) Info(s string, v ...any) {
	logger.infoLogger.Printf(s+"\n", v...)
}

func (logger *CustomLog) Warn(s string, v ...any) {
	logger.infoLogger.Printf(s+"\n", v...)
}

func (logger *CustomLog) Error(s error, v ...any) {
	logger.errorLogger.Printf(s.Error()+"\n", v...)
}
