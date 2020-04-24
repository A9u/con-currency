package service

import (
	"con-currency/config"
	"con-currency/db"
	"con-currency/model"
	"con-currency/xeservice"
	"database/sql"
	"strconv"

	"sync"

	logger "github.com/sirupsen/logrus"
)

//StartProcess start the process of fetching currency exchange rates and insert it into database
func StartProcess() {

	var w sync.WaitGroup
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
		logger.WithField("err", err.Error()).Error("Cannot initialize database")
		return
	}

	//Spawning goroutines for processing each currency
	for _, currency := range currencies {
		w.Add(1)
		go apiToDB(currency, dbInstance, &w)
	}

	w.Wait()

}

func apiToDB(currency string, dbInstance *sql.DB, wg *sync.WaitGroup) {

	defer wg.Done()

	resp, err := xeservice.GetExRateFromAPI(currency)
	if err != nil {
		logger.WithField("err", err).Error("API issue found")
		return
	}

	query, val := queryBuilder(resp)

	result, err := db.FireQuery(query, val, dbInstance)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Query execution failed")
		return
	}

	rowCnt, err := result.RowsAffected()
	if err != nil {
		logger.WithField("err", err.Error()).Error("Update DB failed")
		return
	}
	logger.WithField("affected rows", rowCnt).Info("Update DB successful")
	return
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
