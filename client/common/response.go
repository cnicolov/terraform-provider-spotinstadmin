package common

import "encoding/json"

type Response struct {
	Request struct {
		ID string `json:"id"`
	} `json:"request"`

	Response struct {
		Errors []ResponseError   `json:"errors"`
		Items  []json.RawMessage `json:"items"`
	} `json:"response"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
}
