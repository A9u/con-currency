package db

import (
	"database/sql"
	"fmt"
	"strconv"

	"con-currency/config"
	"con-currency/model"

	_ "github.com/lib/pq" //pgDriver
	logger "github.com/sirupsen/logrus"
)

const (
	driverName = "postgres"
)

// Init Initialize DB
func Init() (db *sql.DB, err error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.GetConfigString("db_config.host"),
		config.GetConfigString("db_config.port"),
		config.GetConfigString("db_config.user"),
		config.GetConfigString("db_config.password"),
		config.GetConfigString("db_config.dbname"),
		config.GetConfigString("db_config.sslmode"))

	db, err = sql.Open(driverName, conStr)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Cannot open connection")
		return
	}

	logger.WithField("conn string", conStr).Info("DB connected successfully")

	err = createTable(db)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Cannot create table")
		return
	}

	return
}

// UpdateCurrencies insert/update currencies into database
func UpdateCurrencies(xeResp model.XEcurrency, db *sql.DB) (rowCnt int64, err error) {

	query, val := queryBuilder(xeResp)

	result, err := db.Exec(query, val...)
	if err != nil {
		fmt.Println("query→", query, val)
		logger.WithField("err", err.Error()).Error("Query failed")
		return 0, err
	}

	rowCnt, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowCnt, nil
}

func createTable(db *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS public.exchange_rates(
	      from_currency character varying(3) COLLATE pg_catalog."default",
	      to_currency character varying(3) COLLATE pg_catalog."default",
	      rate double precision,
	      created_at timestamp with time zone,
	      updated_at timestamp with time zone,
          CONSTRAINT unq UNIQUE (from_currency, to_currency))`

	_, err := db.Exec(q)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Create table failed")
		return err
	}

	return nil
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
