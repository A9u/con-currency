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

	exchangeRater := xeservice.New() // will ret interface

	storer, err := db.Init() // will ret interface
	if err != nil {
		logger.WithField("err", err.Error()).Error("Cannot connect database")
		return
	}

	// close database connection
	defer storer.Close()
	if err = storer.CreateTableIfMissing(); err != nil {
		logger.WithField("err", err.Error()).Error("Cannot create table")
	}

	currencies := config.GetStringSlice("currency_list")

	// Starting the process
	service.StartProcess(currencies, exchangeRater, storer)

	elapsed := time.Since(start)
	logger.WithField("info:", elapsed).Info("Execution time")
}
