package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/nsqsink/washtub/internal/models"
	"github.com/nsqsink/washtub/internal/repositories"
	"github.com/nsqsink/washtub/internal/usecases"
	"github.com/nsqsink/washtub/pkg/sock"
)

type Server struct {
	port   int
	memDB  *memdb.MemDB
	socHub *sock.Hub

	// Repos
	workerStore  models.WorkerStore
	messageStore models.MessageStore

	// Usecases
	workerUsecase  models.WorkerUsecase
	messageUsecase models.MessageUsecase
}

func NewServer(port int, memDB *memdb.MemDB, socHub *sock.Hub) *http.Server {
	server := &Server{
		port:   port,
		memDB:  memDB,
		socHub: socHub,
	}
	server.init()

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", server.port),
		Handler:      server.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) init() {
	// Repos
	s.workerStore = repositories.NewWorkerStore(s.memDB)
	s.messageStore = repositories.NewMessageStore(s.memDB)

	// Usecases
	s.workerUsecase = usecases.NewWorkerUsecase(s.workerStore, s.socHub)
	s.messageUsecase = usecases.NewMessageUsecase(s.messageStore)
}
