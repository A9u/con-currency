package db

import (
	"database/sql"
	"fmt"

	"con-currency/config"

	_ "github.com/lib/pq" //pgDriver
	logger "github.com/sirupsen/logrus"
)

const (
	driverName = "postgres"
)

// Init Initialize DB
func Init() (db *sql.DB, err error) {

	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.GetConfig("db_config.host"),
		config.GetConfig("db_config.port"),
		config.GetConfig("db_config.user"),
		config.GetConfig("db_config.password"),
		config.GetConfig("db_config.dbname"),
		config.GetConfig("db_config.sslmode"))

	db, err = sql.Open(driverName, conStr)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Cannot open connection")
		return
	}
	logger.WithField("conn string", conStr).Info("DB connected successfully")
	return
}

// FireQuery execute db query
func FireQuery(query string, val []interface{}, db *sql.DB) (result sql.Result, err error) {

	result, err = db.Exec(query, val...)
	if err != nil {
		fmt.Println("query→", query, val)
		logger.WithField("err", err.Error()).Error("Query failed")
		return result, err
	}

	return result, nil
}

// CreateTable if not exist
func CreateTable(db *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS public.exchange_rates(
	       from_currency character varying(3) COLLATE pg_catalog."default",
	       to_currency character varying(3) COLLATE pg_catalog."default",
	       rate double precision,
	       created_at timestamp with time zone,
	   	updated_at timestamp with time zone,
	   	CONSTRAINT unq UNIQUE (from_currency, to_currency)
	   )`
	_, err := db.Exec(q)
	if err != nil {

		logger.WithField("err", err.Error()).Error("Create table failed")
		return err
	}

	return nil

}
