package service

import (
	"con-currency/db"
	"con-currency/model"
	"con-currency/xeservice"

	logger "github.com/sirupsen/logrus"
)

// type XEServiceMock struct {
// 	URL      string
// 	Username string
// 	Password string
// }

// func (xeService XEServiceMock) GetExchangeRate(currency string) (xeResp model.XEcurrency, err error) {
// 	return
// }

//StartProcess start the process of fetching currency exchange rates and insert it into database
func StartProcess(currencies []string, exchangeRater xeservice.ExchangeRater, storer db.Storer) {
	var rowsAffected int64

	//xe := XEServiceMock{}
	// creating channel for sending jobs
	jobs := make(chan string, len(currencies))

	// creating channel for recieving errors and response
	results := make(chan model.Results, len(currencies))

	// Creating workers
	for w := 0; w <= 10; w++ {
		go processCurrencies(exchangeRater, storer, jobs, results)
	}

	// sending jobs
	for _, currency := range currencies {
		jobs <- currency
	}

	close(jobs)

	// recieving results
	for i := 0; i < len(currencies); i++ {
		res := <-results
		if res.Err != nil {
			logger.WithField("err", res.Err.Error()).Error("Exit")
			return
		}

		rowsAffected += res.RowsAffected
	}

	logger.WithField("rows affected", rowsAffected).Info("Job successfull")
}

// func processCurrencies(xeService xeservice.GetExchangeRater, dbInstance *sql.DB, jobs <-chan string, results chan<- model.Results) {
func processCurrencies(exchangeRater xeservice.ExchangeRater, storer db.Storer, jobs <-chan string, results chan<- model.Results) {

	for currency := range jobs {
		rowCnt, err := processCurrency(currency, exchangeRater, storer)
		if err != nil {
			results <- model.Results{
				RowsAffected: 0,
				Err:          err,
			}
			return
		}

		results <- model.Results{
			RowsAffected: rowCnt,
			Err:          nil,
		}

	}
}

// func processCurrency(currency string, xeService xeservice.GetExchangeRater, dbInstance *sql.DB) (rowCnt int64, err error) {
func processCurrency(currency string, exchangeRater xeservice.ExchangeRater, storer db.Storer) (rowCnt int64, err error) {
	xeResp, err := exchangeRater.GetExchangeRate(currency)
	if err != nil {
		return
	}

	rowCnt, err = storer.UpsertCurrencies(xeResp)
	if err != nil {
		return
	}

	return
}
