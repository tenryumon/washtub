package usecases

import (
	"context"

	"github.com/nsqsink/washtub/internal/models"
)

type workerUsecase struct {
	workerRepo models.WorkerRepository
}

func NewWorkerUsecase(repo models.WorkerRepository) models.WorkerUsecase {
	return &workerUsecase{
		workerRepo: repo,
	}
}

// Fetch implements models.WorkerUsecase.
func (w *workerUsecase) Fetch(ctx context.Context, request models.FetchRequest) (res []models.Worker, err error) {
	res, err = w.workerRepo.Fetch(ctx, request)
	if err != nil {
		return nil, err
	}
	return
}

// GetByID implements models.WorkerUsecase.
func (w *workerUsecase) GetByID(ctx context.Context, id string) (res models.Worker, err error) {
	res, err = w.workerRepo.GetByID(ctx, id)
	return
}

// Pulse implements models.WorkerUsecase.
func (w *workerUsecase) Pulse(ctx context.Context, worker models.Worker) (err error) {
	err = w.workerRepo.Store(ctx, worker)
	return
}
