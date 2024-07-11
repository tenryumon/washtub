package handlers

import (
	"net/http"
)

type HealthCheck struct{}

func NewHealthcheckHandler() *HealthCheck {
	return &HealthCheck{}
}

func (h *HealthCheck) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HealthCheck) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
