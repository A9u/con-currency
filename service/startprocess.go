package service

import (
	"con-currency/config"
	"con-currency/db"
	"con-currency/model"
	"con-currency/xeservice"
	"database/sql"
	"runtime"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

//StartProcess start the process of fetching currency exchange rates and insert it into database
func StartProcess() {
	var rowsAffected int64

	var currencies = config.GetStringSlice("currency_list")

	logger.WithField("currencies", currencies).Info("Currencies initialized")

	//Initialize database
	dbInstance, err := db.Init()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Cannot initialize database")
		return
	}
	defer dbInstance.Close()

	//Create table if not exist
	err = db.CreateTable(dbInstance)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Cannot create table")
		return
	}

	// creating channel for recieving errors and response
	results := make(chan model.Result, len(currencies))

	// creating channel for sending jobs
	jobs := make(chan string, len(currencies))

	// Creating workers
	for w := 0; w <= runtime.NumCPU()-1; w++ {
		go apiToDB(dbInstance, jobs, results)

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

func apiToDB(dbInstance *sql.DB, jobs <-chan string, results chan<- model.Results) {
	for currency := range jobs {
		xeResp, err := xeservice.GetExRateFromAPI(currency)
		if err != nil {
			results <- model.Results{
				0,
				err,
			}
			return
		}

		query, val := queryBuilder(xeResp)

		dbResp, err := db.FireQuery(query, val, dbInstance)

		if err != nil {
			results <- model.Results{
				0,
				err,
			}
			return
		}

		rowCnt, err := dbResp.RowsAffected()
		if err != nil {
			results <- model.Results{
				0,
				err,
			}
			return
		}

		results <- model.Results{
			rowCnt,
			nil,
		}

		return
	}

}

func queryBuilder(resp model.XEcurrency) (string, []interface{}) {
	values := []interface{}{}
	query := `INSERT INTO exchange_rates (from_currency,to_currency,rate,created_at,updated_at) values `

	for i, r := range resp.To {
		//appending keys
		values = append(values, resp.From, r.Quotecurrency, r.Mid, resp.Timestamp, resp.Timestamp)

		numFields := 5
		n := i * numFields

		//appending $1, $2, ...
		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$` + strconv.Itoa(n+j+1) + `,`
		}
		query = query[:len(query)-1] + `),`

	}

	query = query[:len(query)-1]
	query += `ON CONFLICT ON CONSTRAINT unq
		DO UPDATE SET rate =excluded.rate,updated_at = excluded.updated_at where exchange_rates.rate is distinct from excluded.rate`

	return query, values

}
