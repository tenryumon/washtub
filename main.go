package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/nsqsink/washtub/internal/handlers"
	"github.com/nsqsink/washtub/internal/repositories"
	"github.com/nsqsink/washtub/internal/usecases"
	"golang.org/x/exp/slog"
)

type Configuration struct {
	Database DatabaseConfig
	Http     HttpConfig
}

type DatabaseConfig struct {
	Driver     string
	Connection string
}
type HttpConfig struct {
	BaseDomain string
	Port       string
}

func main() {
	// Get Flag for service
	var configFile string
	flag.StringVar(&configFile, "config", "files/configuration/config.yml", "Configuration File Location")
	flag.Parse()

	r := chi.NewRouter()

	r.Get("/ping", ping)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Worker
	workerRepo := repositories.NewWorkerRepository()
	workerUsecase := usecases.NewWorkerUsecase(workerRepo)
	handlers.NewWorkerHandler(r, workerUsecase)

	srv := &http.Server{
		Addr:    "0.0.0.0:9000",
		Handler: r,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			slog.Error("HTTP server Shutdown", err)
		}
		close(idleConnsClosed)
	}()

	slog.Info("Start HTTP Server Port", 9000)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("HTTP server ListenAndServe", err)
	}

	<-idleConnsClosed
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
