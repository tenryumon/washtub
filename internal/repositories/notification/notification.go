package notification

import (
	"context"
	"errors"
	ht "html/template"
	tt "text/template"

	"github.com/nyelonong/boilerplate-go/core/log"
	"github.com/nyelonong/boilerplate-go/core/mailer"
	"github.com/nyelonong/boilerplate-go/core/meta"
	"github.com/nyelonong/boilerplate-go/core/qiscus"
	"github.com/nyelonong/boilerplate-go/internal/entities"
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	SmsTemplateLoc      string
	WhatsappTemplateLoc string
	EmailTemplateLoc    string
	Mailer              *mailer.Mailer
	Whatsapp            *meta.Whatsapp
	Qiscus              *qiscus.Qiscus
	emailSenderMap      map[string]entities.EmailSender
}

type NotificationObject struct {
	smsTemplates      *tt.Template
	whatsappTemplates *tt.Template
	emailTemplates    *ht.Template
	mailer            *mailer.Mailer
	whatsapp          *meta.Whatsapp
	qiscus            *qiscus.Qiscus
	emailSenderMap    map[string]entities.EmailSender
}

func (c *Configuration) AddEmailSender(key, name, email string) error {
	if !entities.IsValidEmailSender(key) {
		return errors.New("invalid email sender")
	}
	if c.emailSenderMap == nil {
		c.emailSenderMap = map[string]entities.EmailSender{}
	}
	c.emailSenderMap[key] = entities.EmailSender{
		Name:  name,
		Email: email,
	}
	return nil
}

func New(config Configuration) interfaces.Notification {
	return &NotificationObject{
		smsTemplates:      tt.Must(tt.ParseGlob(config.SmsTemplateLoc + "/*.txt")),
		whatsappTemplates: tt.Must(tt.ParseGlob(config.WhatsappTemplateLoc + "/*.txt")),
		emailTemplates:    ht.Must(ht.ParseGlob(config.EmailTemplateLoc + "/*.html")),
		mailer:            config.Mailer,
		whatsapp:          config.Whatsapp,
		qiscus:            config.Qiscus,
		emailSenderMap:    config.emailSenderMap,
	}
}

func (no *NotificationObject) Send(ctx context.Context, channel entities.Channel, metadata map[string]interface{}) error {
	err := no.sendByType(channel, metadata)
	if err != nil {
		log.Errorf("Failed to do Send via %s to %s because %s", channel.Type, channel.Destination, err)
		// Bypass Error for now
	}

	return nil
}

func (no *NotificationObject) SendMulti(ctx context.Context, channels []entities.Channel, metadata map[string]interface{}) error {
	for _, channel := range channels {
		err := no.sendByType(channel, metadata)
		if err != nil {
			log.Errorf("Failed to do Send via %s to %s because %s", channel.Type, channel.Destination, err)
			// Bypass Error for now
		}
	}

	return nil
}

func (no *NotificationObject) sendByType(channel entities.Channel, metadata map[string]interface{}) error {
	var err error

	switch channel.Type {
	case entities.TypeEmail:
		err = no.sendEmail(channel, metadata)
	case entities.TypeSMS:
		err = no.sendSMS(channel.Destination, channel.Template, metadata)
	case entities.TypeWhatsapp:
		err = no.sendWhatsapp(channel.Destination, channel.Template, metadata)
	case entities.TypeQiscusWhatsapp:
		err = no.sendQiscusWhatsapp(channel.Destination, metadata)
	}

	return err
}
