package models

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/go-memdb"
	httpLib "github.com/nsqsink/washtub/pkg/http"
)

type (
	Message struct {
		ID        string    `json:"id"`
		WorkerID  string    `json:"worker_id"`
		Body      string    `json:"body"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	MessageUsecase interface {
		Store(ctx context.Context, message Message) (Message, error)
		Fetch(ctx context.Context, request httpLib.FetchRequest, workerId string) ([]Message, error)
		GetByID(ctx context.Context, id string) (Message, error)
	}
	MessageStore interface {
		Store(ctx context.Context, message Message) (Message, error)
		Fetch(ctx context.Context, request httpLib.FetchRequest, workerId string) (res []Message, err error)
		GetByID(ctx context.Context, id string) (Message, error)
		Stash(ctx context.Context, message Message) error
	}
)

func (m *Message) SetID() {
	raw := fmt.Sprintf("%s-%s", m.WorkerID, m.Body)
	m.ID = base64.StdEncoding.EncodeToString([]byte(raw))
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
