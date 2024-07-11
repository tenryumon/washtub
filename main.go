package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/hashicorp/go-memdb"
	"github.com/nsqsink/washtub/internal/models"
	"github.com/nsqsink/washtub/internal/server"
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

	// Server
	initServer(memDB, hub)
}

func initServer(memDB *memdb.MemDB, hub *sock.Hub) {
	srv := server.NewServer(9000, memDB, hub)

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
