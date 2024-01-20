package http

import (
	"net/http"
	"strings"
)

type FetchRequest struct {
	Search string `json:"search"`
	Sort   Sort   `json:"sort"`
}

type Sort struct {
	Key       string `json:"key"`
	Direction string `json:"direction"`
}

func BuildFetchRequest(r *http.Request) (req FetchRequest) {
	req.Search = r.URL.Query().Get("search")

	if sort := r.URL.Query().Get("sort"); sort != "" {
		str := strings.Split(sort, ",")
		if len(str) > 1 {
			req.Sort = Sort{
				Key:       str[0],
				Direction: str[1],
			}
		}
	}

	return
}
