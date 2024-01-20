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
	Worker struct {
		ID        string    `json:"id"`
		ChannelID string    `json:"channel_id"`
		Address   string    `json:"address"`
		Topic     string    `json:"topic"`
		SinkType  string    `json:"sink_type"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	WorkerUsecase interface {
		Pulse(ctx context.Context, worker Worker) (Worker, error)
		Fetch(ctx context.Context, request httpLib.FetchRequest) ([]Worker, error)
		GetByID(ctx context.Context, id string) (Worker, error)
	}
	WorkerStore interface {
		Store(ctx context.Context, worker Worker) (Worker, error)
		Fetch(ctx context.Context, request httpLib.FetchRequest) (res []Worker, err error)
		GetByID(ctx context.Context, id string) (Worker, error)
		Stash(ctx context.Context, worker Worker) error
	}
)

const (
	STATUS_ACTIVE   = "active"
	STATUS_INACTIVE = "inactive"
)

func (w *Worker) SetID() {
	raw := fmt.Sprintf("%s-%s-%s-%s", w.ChannelID, w.Address, w.Topic, w.SinkType)
	w.ID = base64.StdEncoding.EncodeToString([]byte(raw))
}

var WorkerSchema = &memdb.TableSchema{
	Name: "worker",
	Indexes: map[string]*memdb.IndexSchema{
		"id": {
			Name:    "id",
			Unique:  true,
			Indexer: &memdb.StringFieldIndex{Field: "ID"},
		},
		"channelid": {
			Name:    "channelid",
			Unique:  true,
			Indexer: &memdb.StringFieldIndex{Field: "ChannelID"},
		},
		"sinktype": {
			Name:    "sinktype",
			Unique:  false,
			Indexer: &memdb.StringFieldIndex{Field: "SinkType"},
		},
		"status": {
			Name:    "status",
			Unique:  false,
			Indexer: &memdb.StringFieldIndex{Field: "Status"},
		},
	},
}
