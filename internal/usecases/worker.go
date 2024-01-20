package usecases

import (
	"context"
	"fmt"
	"reflect"

	"github.com/nsqsink/washtub/internal/models"
	httpLib "github.com/nsqsink/washtub/pkg/http"
	"github.com/nsqsink/washtub/pkg/sock"
)

type WorkerUsecase struct {
	WorkerStore models.WorkerStore
	bucket      *sock.Hub
}

func NewWorkerUsecase(store models.WorkerStore, hub *sock.Hub) models.WorkerUsecase {
	return &WorkerUsecase{
		WorkerStore: store,
		bucket:      hub,
	}
}

// Fetch implements models.WorkerUsecase.
func (w *WorkerUsecase) Fetch(ctx context.Context, request httpLib.FetchRequest) (res []models.Worker, err error) {
	res, err = w.WorkerStore.Fetch(ctx, request)
	if err != nil {
		return nil, err
	}
	return
}

// GetByID implements models.WorkerUsecase.
func (w *WorkerUsecase) GetByID(ctx context.Context, id string) (res models.Worker, err error) {
	res, err = w.WorkerStore.GetByID(ctx, id)
	return
}

// Pulse implements models.WorkerUsecase.
func (w *WorkerUsecase) Pulse(ctx context.Context, worker models.Worker) (res models.Worker, err error) {
	// SetID before process
	worker.SetID()

	// Check data by ID
	existing, err := w.WorkerStore.GetByID(ctx, worker.ID)
	if err != nil {
		return res, err
	}

	// Store data to memory DB
	res, err = w.WorkerStore.Store(ctx, worker)
	if err != nil {
		return res, err
	}

	// If data is different, broadcast to client
	if (models.Worker{}) == existing || !reflect.DeepEqual(existing, res) {
		w.bucket.Broadcast(fmt.Sprintf("worker:%s", res.ID))
	}
	return
}
