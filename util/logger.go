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

func (logger *Logger) Info(s string, v ...any) {
	logger.infoLogger.Printf(s+"\n", v...)
}

func (logger *Logger) Warn(s string, v ...any) {
	logger.infoLogger.Printf(s+"\n", v...)
}

func (logger *Logger) Error(s string, v ...any) {
	logger.errorLogger.Printf(s+"\n", v...)
}
