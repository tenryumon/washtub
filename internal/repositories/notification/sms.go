package notification

import (
	"bytes"
	"fmt"
)

func (no *NotificationObject) sendSMS(destination string, template string, metadata map[string]interface{}) error {
	var tpl bytes.Buffer

	err := no.smsTemplates.ExecuteTemplate(&tpl, template+".txt", metadata)
	if err != nil {
		return fmt.Errorf("Failed to ExecuteTemplate %s because %s", template, err)
	}
	println(tpl.String())

	// Send SMS Here

	return nil
}
