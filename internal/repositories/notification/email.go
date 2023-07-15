package notification

import (
	"bytes"
	"fmt"

	"github.com/nyelonong/boilerplate-go/core/mailer"
	"github.com/nyelonong/boilerplate-go/internal/entities"
)

func (no *NotificationObject) sendEmail(channel entities.Channel, metadata map[string]interface{}) error {
	var tpl bytes.Buffer

	err := no.emailTemplates.ExecuteTemplate(&tpl, channel.Template+".html", metadata)
	if err != nil {
		return fmt.Errorf("Failed to ExecuteTemplate %s because %s", channel.Template, err)
	}
	if _, ok := no.emailSenderMap[channel.EmailSender]; !ok {
		return fmt.Errorf("Failed to map email sender with name %s", channel.EmailSender)
	}

	// Send Email Here
	emailFrom := no.emailSenderMap[channel.EmailSender]
	_, err = no.mailer.SendEmails(mailer.MailOptions{
		EmailFrom: emailFrom.Email,
		EmailTo:   []string{channel.Destination},
		FromName:  emailFrom.Name,
		Subject:   channel.Subject,
		Content:   tpl.String(),
	})
	if err != nil {
		return err
	}

	return nil
}
