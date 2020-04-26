package model

import "time"

// XEcurrency is model of API response
type XEcurrency struct {
	Terms     string    `json:"terms"`
	Privacy   string    `json:"privacy"`
	From      string    `json:"from"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	To        []struct {
		Quotecurrency string  `json:"quotecurrency"`
		Mid           float64 `json:"mid"`
	} `json:"to"`
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
