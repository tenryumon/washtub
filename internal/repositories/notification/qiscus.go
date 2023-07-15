package notification

import (
	"fmt"
)

func (no *NotificationObject) sendQiscusWhatsapp(destination string, metadata map[string]interface{}) error {
	otp := ""
	if val, ok := metadata["otp"]; ok {
		otp = val.(string)
	}

	// Send Qiscus Template Here
	msgBody := fmt.Sprintf(`
		{
			"to": "%s",
			"type": "template",
			"template": {
				"namespace": "%s",
				"name": "otp",
				"language": {
					"policy": "deterministic",
					"code": "id"
				},
				"components": [
					{
						"type" : "body",
						"parameters": [
							{
								"type": "text",
								"text": "%s"
							}
						]
					}
				]
			}
		}
	`, formatPhone(destination), no.qiscus.GetNamespace(), otp)

	err := no.qiscus.SendWhatsapp(msgBody)
	if err != nil {
		return err
	}
	return nil
}
