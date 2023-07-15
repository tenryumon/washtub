package interfaces

import (
	"context"

	"github.com/nyelonong/boilerplate-go/internal/models"
)

type GeneralUC interface {
	// Upload Image and File
	UploadImage(ctx context.Context, param models.UploadFileReq) models.UploadFileResp
	UploadFile(ctx context.Context, param models.UploadFileReq) models.UploadFileResp

	// Webhook
	WebhookXenditInvoices(ctx context.Context, param models.WebhookXenditInvoicesReq) (resp models.EmptyDataResp)
}
