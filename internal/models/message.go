package models

import (
	"context"
	"time"

	"github.com/hashicorp/go-memdb"
)

type Message struct {
	ID        string    `json:"id"`
	WorkerID  string    `json:"worker_id"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MessageUsecase interface {
	Store(ctx context.Context, message Message) error
	Fetch(ctx context.Context, request FetchRequest, workerId string) ([]Message, error)
	GetByID(ctx context.Context, id string) (Message, error)
}

type MessageStore interface {
	Store(ctx context.Context, message Message) error
	Fetch(ctx context.Context, request FetchRequest, workerId string) (res []Message, err error)
	GetByID(ctx context.Context, id string) (Message, error)
	Stash(ctx context.Context, message Message) error
}

var MessageSchema = &memdb.TableSchema{
	Name: "message",
	Indexes: map[string]*memdb.IndexSchema{
		"id": {
			Name:    "id",
			Unique:  true,
			Indexer: &memdb.StringFieldIndex{Field: "ID"},
		},
		"workerid": {
			Name:    "workerid",
			Unique:  false,
			Indexer: &memdb.StringFieldIndex{Field: "WorkerID"},
		},
		"status": {
			Name:    "status",
			Unique:  false,
			Indexer: &memdb.StringFieldIndex{Field: "Status"},
		},
	},
}
