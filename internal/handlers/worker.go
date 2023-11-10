package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nsqsink/washtub/internal/models"
)

type WorkerHandler struct {
	WorkerUsecase models.WorkerUsecase
}

func NewWorkerHandler(r chi.Router, uc models.WorkerUsecase) {
	handler := &WorkerHandler{
		WorkerUsecase: uc,
	}

	r.Route("/worker", func(r chi.Router) {
		r.Post("/pulse", handler.Pulse)
	})
}

func (h *WorkerHandler) Pulse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		worker models.Worker
	)

	// Get Request Body
	err := json.NewDecoder(r.Body).Decode(&worker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		panic(err)
	}

	// Call Usecase
	err = h.WorkerUsecase.Pulse(ctx, worker)
	if err != nil {
		panic(err)
	}

	// Write response
	render.JSON(w, r, worker)
}
