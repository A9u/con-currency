package service

import (
	"con-currency/db"
	"con-currency/exchangerate"
	logger "github.com/sirupsen/logrus"
)

//StartProcess start the process of fetching currency exchange rates and insert it into database
func StartProcess(currencies []string, converter exchangerate.Converter, storer db.Storer) {
	var rowsAffected int64

	// creating channel for receiving errors and response
	resultChan := make(chan int64, len(currencies))
	currencyChan := make(chan string, len(currencies))

	for i := 0; i <= 11; i++ {
		go processCurrency(converter, storer, currencyChan, resultChan)
	}
	// sending jobs
	for _, currency := range currencies {
		currencyChan <- currency
	}

	close(currencyChan)
	// recieving results
	for i := 0; i < len(currencies); i++ {
		res := <-resultChan
		rowsAffected += res

	}

	logger.WithField("rows affected", rowsAffected).Info("Job successfull")
	close(resultChan)
}

// func processCurrency(currency string, xeService xeservice.GetConverter, dbInstance *sql.DB) (rowCnt int6464, err error) {
func processCurrency(converter exchangerate.Converter, storer db.Storer, currencyChan <-chan string, results chan<- int64) {
	for currency := range currencyChan {
		resp, err := converter.Get(currency)
		if err != nil {
			return
		}

		rowCnt, err := storer.UpsertCurrencies(resp)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Exit")
			return
		}

		results <- rowCnt
	}
}
