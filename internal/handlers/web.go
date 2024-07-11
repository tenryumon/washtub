package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/nsqsink/washtub/internal/models"
	httpLib "github.com/nsqsink/washtub/pkg/http"
)

type WebHandler struct {
	WorkerUsecase  models.WorkerUsecase
	MessageUsecase models.MessageUsecase
}

func NewWebHandler(workerUsecase models.WorkerUsecase, messageUsecase models.MessageUsecase) *WebHandler {
	return &WebHandler{
		WorkerUsecase:  workerUsecase,
		MessageUsecase: messageUsecase,
	}
}

func (h *WebHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get Request Param
	request := httpLib.BuildFetchRequest(r)

	// Call Usecase
	workers, err := h.WorkerUsecase.Fetch(ctx, request)
	if err != nil {
		log.Fatalf("Could not fetch workers. Err: %v", err)
	}

	// Write response
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	err = tmpl.Execute(w, workers)
	if err != nil {
		log.Fatalf("Could not execute template. Err: %v", err)
	}
}
