package general_uc

import (
	"context"
	"encoding/json"

	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/internal/models"
)

func (gu *GeneralUsecase) WebhookXenditInvoices(ctx context.Context, param models.WebhookXenditInvoicesReq) (resp models.EmptyDataResp) {
	// Initialize Usecase Parameter
	resp.InitError()
	ucContext := map[string]interface{}{
		"usecase": "WebhookXenditInvoices",
	}

	_, err := json.Marshal(param.Metadata)
	if err != nil {
		log.ErrorWithField(ucContext, "Failed to do Marshal Metadata because %s", err)
		return
	}

	// Start Usecase Action
	// xenCallback := entities.XenditCallback{
	// 	ID:           uuid.New().String(),
	// 	Payload:      param.Payload,
	// 	Metadata:     metadata,
	// 	PartitionKey: 0,
	// 	Status:       entities.XenditCallbackStatusTodo,
	// 	CreatedBy:    "0",
	// }
	// err = gu.admission.InsertXenditCallback(ctx, xenCallback)
	// if err != nil {
	// 	log.ErrorWithField(ucContext, "Failed to do InsertXenditCallback because %s", err)
	// 	return
	// }
	resp.Success()
	return
}
