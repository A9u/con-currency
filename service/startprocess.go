package service

import (
	"con-currency/db"
	"con-currency/exchangerate"
	"con-currency/model"

	logger "github.com/sirupsen/logrus"
)

// type XEServiceMock struct {
// 	URL      string
// 	Username string
// 	Password string
// }

// func (xeService XEServiceMock) GetConverte(currency string) (xeResp model.XEcurrency, err error) {
// 	return
// }

//StartProcess start the process of fetching currency exchange rates and insert it into database
func StartProcess(currencies []string, converter exchangerate.Converter, storer db.Storer) {
	var rowsAffected int64

	//xe := XEServiceMock{}
	// creating channel for sending jobs
	jobs := make(chan string, len(currencies))

	// creating channel for recieving errors and response
	results := make(chan model.Results, len(currencies))

	// Creating workers
	for w := 0; w <= 10; w++ {
		go processCurrencies(converter, storer, jobs, results)
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

// func processCurrencies(xeService xeservice.GetConverter, dbInstance *sql.DB, jobs <-chan string, results chan<- model.Results) {
func processCurrencies(converter exchangerate.Converter, storer db.Storer, jobs <-chan string, results chan<- model.Results) {

	for currency := range jobs {
		rowCnt, err := processCurrency(currency, converter, storer)
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

// func processCurrency(currency string, xeService xeservice.GetConverter, dbInstance *sql.DB) (rowCnt int64, err error) {
func processCurrency(currency string, converter exchangerate.Converter, storer db.Storer) (rowCnt int64, err error) {
	resp, err := converter.Get(currency)
	if err != nil {
		return
	}

	rowCnt, err = storer.UpsertCurrencies(resp)
	if err != nil {
		return
	}

	return
}
