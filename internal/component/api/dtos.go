package api

import "time"

type DataResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Rate struct {
	CurrencyCode string    `json:"currency_code"`
	Amount       string    `json:"amount"`
	Date         time.Time `json:"date"`
}
