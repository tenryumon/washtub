package models

import (
	"context"
	"time"
)

type Worker struct {
	ChannelID string    `json:"channel_id"`
	Address   string    `json:"address"`
	Topic     string    `json:"topic"`
	SinkType  string    `json:"sink_type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkerUsecase interface {
	Pulse(ctx context.Context, worker Worker) error
	Fetch(ctx context.Context, request FetchRequest) ([]Worker, error)
	GetByID(ctx context.Context, id string) (Worker, error)
}

type WorkerRepository interface {
	Store(ctx context.Context, worker Worker) error
	Fetch(ctx context.Context, request FetchRequest) (res []Worker, err error)
	GetByID(ctx context.Context, id string) (Worker, error)
	Stash(ctx context.Context, worker Worker) error
}
