package db

import "con-currency/model"

type Storer interface {
	CreateTableIfMissing() (err error)
	UpsertCurrencies(model.XEcurrency) (rowCnt int64, err error)
	Close() (err error)
}
