package main

import (
	"con-currency/config"
	"con-currency/db"
	"con-currency/service"
	"con-currency/xeservice"
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
	err := config.Init("config")
	if err != nil {
		logger.WithField("error in config file", err.Error()).Error("Exit")
		return
	}

	//Creating new xeService object
	xeService := xeservice.New()

	//Initialize database
	dbInstance, err := db.Init() // will ret interface
	if err != nil {
		logger.WithField("err", err.Error()).Error("Cannot initialize database")
		return
	}

	defer dbInstance.Close()

	currencies := config.GetStringSlice("currency_list")

	// Starting the process
	service.StartProcess(currencies, xeService, dbInstance)

	elapsed := time.Since(start)
	logger.WithField("info:", elapsed).Info("Execution time")
}
