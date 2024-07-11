package handlers

import (
	"net/http"

	"github.com/nsqsink/washtub/pkg/sock"
)

type SocketHandler struct {
	hub *sock.Hub
}

func NewSocketHandler(hub *sock.Hub) *SocketHandler {
	return &SocketHandler{
		hub: hub,
	}
}

func (s *SocketHandler) WebSocket(w http.ResponseWriter, r *http.Request) {
	sock.ServeWS(s.hub, w, r)
}
