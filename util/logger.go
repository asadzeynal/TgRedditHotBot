package util

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
}

func NewLogger() *Logger {
	infoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger := log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)

	logger := Logger{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		warnLogger:  warnLogger,
	}

	return &logger
}

func (logger *Logger) Info(s string) {
	logger.infoLogger.Println(s)
}

func (logger *Logger) Warn(s string) {
	logger.infoLogger.Println(s)
}

func (logger *Logger) Error(s string) {
	logger.errorLogger.Println(s)
}
