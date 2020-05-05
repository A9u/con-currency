package exchangerate

import "con-currency/model"

type Converter interface {
	Get(string) ([]model.CurrencyRate, error)
}
