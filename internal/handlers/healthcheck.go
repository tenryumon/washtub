package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewHealthcheckHandler(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
}
