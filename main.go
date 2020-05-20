package main

import (
	"con-currency/config"
	"con-currency/db"
	"con-currency/exchangerate"
	"con-currency/service"
	logger "github.com/sirupsen/logrus"
	"time"
)

func main() {
	// Initialize configurations
	err := config.Init("config")

	if err != nil {
		logger.WithField("error in config file", err.Error()).Error("Exit")
		return
	}

	interval := config.GetInt("time_interval")
	for {
		start := time.Now()
		fetchCurrency()

		elapsed := time.Since(start)
		logger.WithField("info:", elapsed).Info("Execution time")

		elapsedSeconds := elapsed.Round(time.Second).Seconds()
		remaining := float64(interval) - elapsedSeconds
		logger.WithField("info: ", remaining).Info("Sleeping")

		time.Sleep(time.Duration(remaining) * time.Second)
		logger.Info("Waking up")
	}
}

func fetchCurrency() {
	converter := exchangerate.New() // will ret interface

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

	logger.Info(currencies)
	// Starting the process
	service.StartProcess(currencies, converter, storer)
}
