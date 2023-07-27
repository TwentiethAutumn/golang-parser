package logger

import (
	"log"
	"os"
)

type ConsoleLogger struct {
	info    *log.Logger
	warning *log.Logger
}

func (logger *ConsoleLogger) Info(msg string) {
	if logger.info == nil {
		logger.info = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	}

	logger.info.Println(msg)
}

func (logger *ConsoleLogger) Warning(msg string) {
	if logger.warning == nil {
		logger.warning = log.New(os.Stdout, "WARNING\t", log.Ldate|log.Ltime)
	}

	logger.warning.Println(msg)
}
