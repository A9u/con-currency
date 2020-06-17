package exchangerate

import "con-currency/model"

type Converter interface {
	Get(string, string) ([]model.CurrencyRate, error)
}
