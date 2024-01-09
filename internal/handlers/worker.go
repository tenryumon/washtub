package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nsqsink/washtub/internal/models"
	httpResp "github.com/nsqsink/washtub/pkg/http"
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
		r.Get("/fetch", handler.Fetch)
	})
}

func (h *WorkerHandler) Pulse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		request models.Worker
	)

	// Get Request Body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call Usecase
	worker, err := h.WorkerUsecase.Pulse(ctx, request)

	// Write response
	renderer := httpResp.BuildResponseHTTP(worker, err)
	render.Render(w, r, &renderer)
}

func (h *WorkerHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		request models.FetchRequest
		workers []models.Worker
	)

	// Get Request Body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call Usecase
	workers, err = h.WorkerUsecase.Fetch(ctx, request)

	// Write response
	renderer := httpResp.BuildResponseHTTP(workers, err)
	render.Render(w, r, &renderer)
}
