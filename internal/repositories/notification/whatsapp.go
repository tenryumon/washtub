package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyelonong/boilerplate-go/internal/entities"
)

// Will format phone into "62" prefix phone number
func formatPhone(phone string) string {
	if strings.HasPrefix(phone, "62") {
		return phone
	}
	if strings.HasPrefix(phone, "0") {
		return strings.Replace(phone, "0", "62", 1)
	}
	if strings.HasPrefix(phone, "+62") {
		return strings.Replace(phone, "+62", "62", 1)
	}
	return phone
}

func (no *NotificationObject) sendWhatsapp(destination string, template string, metadata map[string]interface{}) error {
	var tpl bytes.Buffer

	err := no.whatsappTemplates.ExecuteTemplate(&tpl, template+".txt", metadata)
	if err != nil {
		return fmt.Errorf("Failed to ExecuteTemplate %s because %s", template, err)
	}

	// Send Whatsapp Here
	msgBody := entities.WhatsappBody{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               formatPhone(destination),
		Type:             "text",
		Text: entities.WhatsappBodyText{
			PreviewURL: false,
			Body:       tpl.String(),
		},
	}

	msgBodyByte, err := json.Marshal(msgBody)
	if err != nil {
		return fmt.Errorf("Failed sendWhatsapp because %s", err)
	}
	err = no.whatsapp.SendWhatsapp(string(msgBodyByte))
	if err != nil {
		return err
	}
	return nil
}
