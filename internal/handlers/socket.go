package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nsqsink/washtub/pkg/sock"
)

func NewSocketHandler(r chi.Router, hub *sock.Hub) {
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		sock.ServeWS(hub, w, r)
	})
}
