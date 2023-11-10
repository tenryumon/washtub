package repositories

import (
	"context"

	"github.com/nsqsink/washtub/internal/models"
)

type workerRepo struct {
}

func NewWorkerRepository() models.WorkerRepository {
	return &workerRepo{}
}

// Fetch implements models.WorkerRepository.
func (*workerRepo) Fetch(ctx context.Context, request models.FetchRequest) (res []models.Worker, err error) {
	panic("unimplemented")
}

// GetByID implements models.WorkerRepository.
func (*workerRepo) GetByID(ctx context.Context, id string) (models.Worker, error) {
	panic("unimplemented")
}

// Stash implements models.WorkerRepository.
func (*workerRepo) Stash(ctx context.Context, worker models.Worker) error {
	panic("unimplemented")
}

// Store implements models.WorkerRepository.
func (*workerRepo) Store(ctx context.Context, worker models.Worker) (err error) {
	return
}
