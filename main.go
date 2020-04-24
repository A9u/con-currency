package main

import (
	"con-currency/config"
	"con-currency/service"
	"time"

	logger "github.com/sirupsen/logrus"
)

func main() {
	start := time.Now()
	// logger config
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})

	// Initialize configurations
	err := config.InitConfig()
	if err != nil {
		logger.WithField("error in config file", err.Error()).Error("Exit")
		return
	}

	// Starting the process
	service.StartProcess()

	elapsed := time.Since(start)
	logger.WithField("info:", elapsed).Info("Execution time")
}
