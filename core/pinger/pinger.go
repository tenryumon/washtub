package pinger

import (
	"encoding/json"
	"fmt"
	"time"
)

var services []PingService

type Pinger interface {
	Ping() error
}
type PingService struct {
	name   string
	pinger Pinger
}

func AddService(name string, service Pinger) {
	services = append(services, PingService{name: name, pinger: service})
}

type PingResponse struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
	Latency string `json:"latency"`
	Error   string `json:"error"`
}

func CheckHealth() []byte {
	var response []PingResponse
	for _, v := range services {
		println(v.name)
		resp := PingResponse{
			Success: true,
			Name:    v.name,
		}

		start := time.Now()
		err := v.pinger.Ping()
		if err != nil {
			resp.Success = false
			resp.Error = err.Error()
		}
		resp.Latency = fmt.Sprintf("%.3fms", time.Now().Sub(start).Seconds()*1000)

		response = append(response, resp)
	}

	result, _ := json.Marshal(response)
	return result
}
