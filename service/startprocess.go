package service

import (
	"con-currency/db"
	"con-currency/exchangerate"
	"con-currency/model"

	logger "github.com/sirupsen/logrus"
)

//StartProcess start the process of fetching currency exchange rates and insert it into database
func StartProcess(currencies []string, converter exchangerate.Converter, storer db.Storer) {
	var rowsAffected int64

	// creating channel for receiving errors and response
	resultChan := make(chan model.Results)

	// sending jobs
	for _, currency := range currencies {
		go processCurrency(converter, storer, currency, resultChan)
	}

	// recieving results
	for i := 0; i < len(currencies); i++ {
		res := <-resultChan
		if res.Err != nil {
			logger.WithField("err", res.Err.Error()).Error("Exit")
			return
		}

		rowsAffected += res.RowsAffected
	}

	logger.WithField("rows affected", rowsAffected).Info("Job successfull")
	close(resultChan)
}

// func processCurrency(currency string, xeService xeservice.GetConverter, dbInstance *sql.DB) (rowCnt int64, err error) {
func processCurrency(converter exchangerate.Converter, storer db.Storer, currency string, results chan<- model.Results) {
	resp, err := converter.Get(currency)
	if err != nil {
		return
	}

	rowCnt, err := storer.UpsertCurrencies(resp)
	if err != nil {
		return
	}

	results <- model.Results{
		RowsAffected: rowCnt,
		Err:          nil,
	}
}
