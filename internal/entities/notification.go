package entities

const (
	TypeEmail          = "email"
	TypeSMS            = "sms"
	TypeWhatsapp       = "whatsapp"
	TypeQiscusWhatsapp = "qiscus_whatsapp"
	TypePushNotif      = "push_notif"

	EmailSenderNoReply = "no_reply"
	EmailSenderSupport = "support"
)

var (
	emailSenderList = map[string]string{
		EmailSenderNoReply: "no_reply",
		EmailSenderSupport: "support",
	}
)

type Channel struct {
	Type        string
	Destination string
	Template    string
	Subject     string
	EmailSender string
}

type EmailSender struct {
	Name  string
	Email string
}

func IsValidEmailSender(sender string) bool {
	_, ok := emailSenderList[sender]
	return ok
}

func GetChannelSMS(destination, template string) Channel {
	return Channel{Type: TypeSMS, Destination: destination, Template: template, Subject: "", EmailSender: ""}
}

func GetChannelWhatsapp(destination, template string) Channel {
	return Channel{Type: TypeWhatsapp, Destination: destination, Template: template, Subject: "", EmailSender: ""}
}
func GetChannelQiscusWhatsapp(destination string) Channel {
	return Channel{Type: TypeQiscusWhatsapp, Destination: destination, Template: "", Subject: "", EmailSender: ""}
}

func GetChannelEmailNoReply(destination, subject, template string) Channel {
	return Channel{Type: TypeEmail, Destination: destination, Template: template, Subject: subject, EmailSender: EmailSenderNoReply}
}

func GetChannelEmailSupport(destination, subject, template string) Channel {
	return Channel{Type: TypeEmail, Destination: destination, Template: template, Subject: subject, EmailSender: EmailSenderSupport}
}

type WhatsappBody struct {
	MessagingProduct string           `json:"messaging_product"`
	RecipientType    string           `json:"recipient_type"`
	To               string           `json:"to"`
	Type             string           `json:"type"`
	Text             WhatsappBodyText `json:"text"`
}
type WhatsappBodyText struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}
