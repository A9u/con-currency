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

type pgStore struct {
	db *sql.DB
}

// Init Initialize DB
func Init() (storer Storer, err error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.GetString("db_config.host"),
		config.GetString("db_config.port"),
		config.GetString("db_config.user"),
		config.GetString("db_config.password"),
		config.GetString("db_config.dbname"),
		config.GetString("db_config.sslmode"),
	)

	db, err := sql.Open(driverName, conStr)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Cannot open connection")
		return
	}

	logger.WithField("conn string", conStr).Info("DB connected successfully")

	return &pgStore{db}, nil
}

// UpdateCurrencies insert/update currencies into database
func (s *pgStore) UpsertCurrencies(currencyRates []model.CurrencyRate) (rowCnt int64, err error) {
	query, val := s.queryBuilder(currencyRates)

	result, err := s.db.Exec(query, val...)
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

func (s *pgStore) CreateTableIfMissing() error {
	q := `CREATE TABLE IF NOT EXISTS public.exchange_rates(
	      from_currency character varying(3) COLLATE pg_catalog."default",
	      to_currency character varying(3) COLLATE pg_catalog."default",
	      rate double precision,
	      created_at timestamp with time zone,
	      updated_at timestamp with time zone,
          CONSTRAINT unq UNIQUE (from_currency, to_currency))`

	_, err := s.db.Exec(q)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Create table failed")
		return err
	}

	return nil
}

func (s *pgStore) queryBuilder(currencyRates []model.CurrencyRate) (string, []interface{}) {
	values := []interface{}{}
	query := `INSERT INTO exchange_rates (from_currency,to_currency,rate,created_at,updated_at) values `
	numFields := 5
	for i, rate := range currencyRates {
		//appending keys
		values = append(values, rate.From, rate.To, rate.Amount, rate.Timestamp, rate.Timestamp)

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
		      DO UPDATE SET rate = excluded.rate,updated_at = excluded.updated_at where exchange_rates.rate is distinct from excluded.rate`

	return query, values
}

func (s *pgStore) Close() (err error) {
	err = s.db.Close()

	return
}
