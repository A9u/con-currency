package model

import "time"

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

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Result struct {
	message string
	err     error
}
