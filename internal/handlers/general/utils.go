package general

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/core/monitor"
	"github.com/nyelonong/boilerplate-go/internal/models"
)

// Simple Get Body JSON
//
// -------------------------
func getBody(r *http.Request, request interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed Read Body: %s", err)
		return err
	}
	log.Debugf("Request: %s", string(body))

	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Errorf("Failed Umarshal Body: %s", err)
	}
	return err
}

// General Response Payload for Staff Dashboard
//
// -------------------------
type Response struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
}

func writeResponseJSON(w http.ResponseWriter, data interface{}, err models.Errors) {
	errors := map[string]string{}
	for key, value := range err.ValidationErrors {
		errors[key] = value[0]
	}

	response := Response{
		Code:    err.StatusCode,
		Message: err.ErrorMessage,
		Errors:  errors,
		Data:    data,
	}

	result, _ := json.Marshal(response)
	status := err.GetStatusCode()
	log.Debugf("Response: %s", string(result))

	w.WriteHeader(status)
	w.Write(result)
}

func monitorAPI(apiName string, start time.Time, err models.Errors) {
	monitor.Record("api_call", start, map[string]string{"api": apiName, "status": err.StatusCode, "usecase": "general"})
}

func (gv *GeneralView) getTime(ctx context.Context) time.Time {
	return time.Now()
}
