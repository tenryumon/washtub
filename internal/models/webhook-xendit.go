package models

type WebhookXenditInvoicesReq struct {
	Payload  string                 `json:"payload"`
	Metadata map[string]interface{} `json:"metadata"`
}
