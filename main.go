package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-memdb"
	"github.com/nsqsink/washtub/internal/handlers"
	"github.com/nsqsink/washtub/internal/models"
	"github.com/nsqsink/washtub/internal/repositories"
	"github.com/nsqsink/washtub/internal/usecases"
	"github.com/nsqsink/washtub/pkg/inmemdb"
	"github.com/nsqsink/washtub/pkg/sock"
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

	// Socket Hub
	hub := sock.NewHub()
	go hub.Run()

	// In Memory DB
	err, memDB := inmemdb.InitDB(map[string]*memdb.TableSchema{
		"worker":  models.WorkerSchema,
		"message": models.MessageSchema,
	})
	if err != nil {
		slog.Error("Failed to Init Memory DB", err)
	}

	// Router
	router := chi.NewRouter()
	initHandler(hub, memDB, router)

	// Server
	initServer(router)
}

func initHandler(hub *sock.Hub, memDB *memdb.MemDB, router chi.Router) {
	// Healthcheck
	handlers.NewHealthcheckHandler(router)
	// Socket
	handlers.NewSocketHandler(router, hub)
	// Worker
	workerStore := repositories.NewWorkerStore(memDB)
	workerUsecase := usecases.NewWorkerUsecase(workerStore, hub)

	// Message
	messageStore := repositories.NewMessageStore(memDB)
	messageUsecase := usecases.NewMessageUsecase(messageStore)

	handlers.NewWorkerHandler(router, workerUsecase, messageUsecase)
}

func initServer(handler http.Handler) {
	srv := &http.Server{
		Addr:    "0.0.0.0:9000",
		Handler: handler,
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
