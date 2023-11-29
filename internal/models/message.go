package models

import (
	"context"
	"time"
)

type Message struct {
	ID        string
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

type MessageRepository interface {
	Store(ctx context.Context, message Message) error
	Fetch(ctx context.Context, request FetchRequest, workerId string) (res []Message, err error)
	GetByID(ctx context.Context, id string) (Message, error)
	Stash(ctx context.Context, message Message) error
}
