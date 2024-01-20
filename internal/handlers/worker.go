package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/nsqsink/washtub/internal/models"
	httpLib "github.com/nsqsink/washtub/pkg/http"
)

type (
	WorkerHandler struct {
		WorkerUsecase  models.WorkerUsecase
		MessageUsecase models.MessageUsecase
	}
	PulseRequest struct {
		ChannelID string `json:"channel_id"`
		Address   string `json:"address"`
		Topic     string `json:"topic"`
		SinkType  string `json:"sink_type"`
		Status    string `json:"status"`
	}
	MessageRequest struct {
		ChannelID   string `json:"channel_id"`
		Address     string `json:"address"`
		Topic       string `json:"topic"`
		SinkType    string `json:"sink_type"`
		Status      string `json:"status"`
		MessageBody string `json:"body"`
	}
)

func NewWorkerHandler(r chi.Router, workerUsecase models.WorkerUsecase, messageUsecase models.MessageUsecase) {
	handler := &WorkerHandler{
		WorkerUsecase:  workerUsecase,
		MessageUsecase: messageUsecase,
	}

	r.Route("/worker", func(r chi.Router) {
		r.Post("/pulse", handler.Pulse)
		r.Post("/message", handler.Message)
		r.Get("/fetch", handler.Fetch)
		r.Get("/{workerID}", handler.FetchMessages)
	})
}

func (h *WorkerHandler) Pulse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		request PulseRequest
	)

	// Get Request Body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call Usecase
	worker, err := h.WorkerUsecase.Pulse(ctx, models.Worker{
		ChannelID: request.ChannelID,
		Address:   request.Address,
		Topic:     request.Topic,
		SinkType:  request.SinkType,
		Status:    request.Status,
	})

	// Write response
	renderer := httpLib.BuildResponseHTTP(worker, err)
	render.Render(w, r, &renderer)
}

func (h *WorkerHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get Request Param
	request := httpLib.BuildFetchRequest(r)

	// Call Usecase
	workers, err := h.WorkerUsecase.Fetch(ctx, request)

	// Write response
	renderer := httpLib.BuildResponseHTTP(workers, err)
	render.Render(w, r, &renderer)
}

func (h *WorkerHandler) Message(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		request MessageRequest
	)

	// Get Request Body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call Usecases
	worker, err := h.WorkerUsecase.Pulse(ctx, models.Worker{
		ChannelID: request.ChannelID,
		Address:   request.Address,
		Topic:     request.Topic,
		SinkType:  request.SinkType,
		Status:    request.Status,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := h.MessageUsecase.Store(ctx, models.Message{
		WorkerID: worker.ID,
		Body:     request.MessageBody,
		Status:   request.Status,
	})

	// Write response
	renderer := httpLib.BuildResponseHTTP(message, err)
	render.Render(w, r, &renderer)
}

func (h *WorkerHandler) FetchMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get Request Param
	request := httpLib.BuildFetchRequest(r)
	workerID := chi.URLParam(r, "workerID")

	// Call Usecase
	messages, err := h.MessageUsecase.Fetch(ctx, request, workerID)

	// Write response
	renderer := httpLib.BuildResponseHTTP(messages, err)
	render.Render(w, r, &renderer)
}
