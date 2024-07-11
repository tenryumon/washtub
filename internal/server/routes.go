package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nsqsink/washtub/internal/handlers"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Handle web static assets
	fs := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Web
	webHandler := handlers.NewWebHandler(s.workerUsecase, s.messageUsecase)
	r.Get("/", webHandler.Index)

	// Common
	healthCheckHandler := handlers.NewHealthcheckHandler()
	r.Get("/health", healthCheckHandler.HealthCheck)
	r.Get("/ping", healthCheckHandler.Ping)

	// Web Socket
	socketHandler := handlers.NewSocketHandler(s.socHub)
	r.Get("/ws", socketHandler.WebSocket)

	// Worker
	workerHandler := handlers.NewWorkerHandler(s.workerUsecase, s.messageUsecase)
	r.Route("/worker", func(r chi.Router) {
		r.Post("/pulse", workerHandler.Pulse)
		r.Post("/message", workerHandler.Message)
		r.Get("/fetch", workerHandler.Fetch)
		r.Get("/{workerID}", workerHandler.FetchMessages)
	})

	return r
}
