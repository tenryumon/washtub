package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

type ResponseHTTP struct {
	ResultStatus ResultStatus `json:"result_status"`
	Data         interface{}  `json:"data"`
}

type ResultStatus struct {
	HttpStatus int      `json:"-"`
	Code       string   `json:"code"`
	Reason     string   `json:"reason"`
	Message    []string `json:"message"`
}

func (rd *ResponseHTTP) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, rd.ResultStatus.HttpStatus)
	return nil
}

func BuildResponseHTTP(result interface{}, err error) ResponseHTTP {
	if err == nil {
		return ResponseHTTP{
			ResultStatus: ResultStatus{
				HttpStatus: http.StatusOK,
				Code:       strconv.Itoa(http.StatusOK),
				Reason:     "OK",
				Message:    []string{"Success"},
			},
			Data: result,
		}
	}
	return ResponseHTTP{
		ResultStatus: ResultStatus{
			HttpStatus: http.StatusInternalServerError,
			Code:       "500",
			Reason:     err.Error(),
			Message:    []string{err.Error()},
		},
		Data: result,
	}
}
