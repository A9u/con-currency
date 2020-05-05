package model

import "time"

// Currency is model of API response
type CurrencyRate struct {
	From      string
	To        string
	Amount    float64
	Timestamp time.Time
}

// ErrorResponse is model of API error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Results is model of the gouroutine responses
type Results struct {
	RowsAffected int64
	Err          error
}
