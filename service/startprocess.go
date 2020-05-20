package service

import (
	"con-currency/config"
	"con-currency/db"
	"con-currency/exchangerate"
	"con-currency/model"
	logger "github.com/sirupsen/logrus"
	"strings"
)

//StartProcess start the process of fetching currency exchange rates and insert it into database
func StartProcess(currencies []string, converter exchangerate.Converter, storer db.Storer) {
	var rowsAffected int64

	// creating channel for receiving errors and response
	resultChan := make(chan model.Result, len(currencies))
	currencyChan := make(chan string, len(currencies))

	for i := 0; i <= 11; i++ {
		go processCurrency(converter, storer, currencyChan, resultChan)
	}

	errors := make(map[error]struct{})
	// sending jobs
	for _, currency := range currencies {
		currencyChan <- currency
	}

	close(currencyChan)
	// recieving results
	for i := 0; i < len(currencies); i++ {
		res := <-resultChan
		rowsAffected += res.RowsAffected
		if res.Err != nil {
			errors[res.Err] = struct{}{}
		}
	}

	logger.WithField("rows affected", rowsAffected).Info("Job successfull")

	if len(errors) > 0 {
		logger.Info("send mail here")
		notify(errors)
	}
	close(resultChan)
}

// func processCurrency(currency string, xeService xeservice.GetConverter, dbInstance *sql.DB) (rowCnt int6464, err error) {
func processCurrency(converter exchangerate.Converter, storer db.Storer, currencyChan <-chan string, results chan<- model.Result) {
	for currency := range currencyChan {
		resp, err := converter.Get(currency)
		if err != nil {
			results <- model.Result{0, err}
			return
		}

		rowCnt, err := storer.UpsertCurrencies(resp)
		logger.Info(err)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Exit")
			results <- model.Result{0, err}
			return
		}

		results <- model.Result{rowCnt, nil}
	}
}

func notify(errors map[error]struct{}) {
	var errorMsg string
	for k, _ := range errors {
		if !strings.Contains(errorMsg, k.Error()) {
			errorMsg += k.Error() + "\n"
		}
	}

	body := `Hi
 We have failed to fetch the exchange rates due to the following reasons(s):-
`
	body += errorMsg

	mailer := NewMailer()
	err := mailer.Send(config.GetStringSlice("mail_recipients"), config.GetString("mail_sender"), "Failed to fetch rates", body)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Exit")
	}
}
