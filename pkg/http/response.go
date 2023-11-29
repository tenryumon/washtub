package http

import (
	"net/http"
	"strconv"
)

type ResponseHTTP struct {
	ResultStatus ResultStatus `json:"result_status"`
	Data         interface{}  `json:"data"`
}

type ResultStatus struct {
	Code    string   `json:"code"`
	Reason  string   `json:"reason"`
	Message []string `json:"message"`
}

func BuildResponseHTTP(result interface{}, err error) ResponseHTTP {
	if err == nil {
		return ResponseHTTP{
			ResultStatus: ResultStatus{
				Code:    strconv.Itoa(http.StatusOK),
				Reason:  "OK",
				Message: []string{"Success"},
			},
			Data: result,
		}
	}
	return ResponseHTTP{
		ResultStatus: ResultStatus{
			Code:    "500",
			Reason:  "Reason",
			Message: []string{err.Error()},
		},
		Data: result,
	}
}
