package service

import (
	"con-currency/config"
	"con-currency/db"
	"con-currency/exchangerate"
	"con-currency/model"
	logger "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var lastMailSent time.Time

//StartProcess start the process of fetching currency exchange rates and insert it into database
func StartProcess(currencies []string, converter exchangerate.Converter, storer db.Storer) {

	var rowsAffected int64

	// creating channel for receiving errors and response
	resultChan := make(chan model.Result, len(currencies))
	currencyChan := make(chan string, len(currencies))

	toCurrencies := strings.Join(currencies, ",")
	for i := 0; i <= 20; i++ {
		go processCurrencies(converter, storer, currencyChan, resultChan, toCurrencies)
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
	logger.WithField("count of errors occurred", len(errors)).Info("Error")
	if len(errors) > 0 && canSendMail() {
		logger.Info("send mail here")
		notify(errors)
	}
	close(resultChan)
}

// func processCurrency(currency string, xeService xeservice.GetConverter, dbInstance *sql.DB) (rowCnt int6464, err error) {
func processCurrencies(converter exchangerate.Converter, storer db.Storer, currencyChan <-chan string, results chan<- model.Result, toCurrencies string) {
	for currency := range currencyChan {
		rowCnt, err := processCurrency(converter, storer, currency, toCurrencies)
		results <- model.Result{RowsAffected: rowCnt, Err: err}
	}
}

func processCurrency(converter exchangerate.Converter, storer db.Storer, fromCurrency, toCurrencies string) (rowCnt int64, err error) {
	resp, err := converter.Get(fromCurrency, toCurrencies)
	if err != nil {
		return
	}
	rowCnt, err = storer.UpsertCurrencies(resp)
	return
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

	subject := "Failed to fetch rates [" + config.GetString("environment") + "]"
	mailer := NewMailer()
	err := mailer.Send(config.GetStringSlice("mail_recipients"), config.GetString("mail_sender"), subject, body)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Exit")
	} else {
		lastMailSent = time.Now()
	}
}

func canSendMail() bool {
	elapsed := time.Since(lastMailSent)
	elapsedHours := elapsed.Round(time.Hour).Hours()
	mailInterval := config.GetInt("mail_interval")
	return elapsedHours > float64(mailInterval)

}
