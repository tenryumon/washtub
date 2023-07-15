package mailer

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math"
	"math/big"
	"os"
	"time"

	gomail "gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
}

// MailOptions hold options for sending mail
type MailOptions struct {
	EmailFrom    string
	EmailTo      []string
	EmailCc      []string
	EmailBcc     []string
	EmailReplyTo []string
	Subject      string
	Content      string
	FromName     string
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	SkipSSL  bool
}

func New(mailerConfig Config) *Mailer {
	dialer := gomail.NewDialer(mailerConfig.Host, mailerConfig.Port, mailerConfig.Username, mailerConfig.Password)
	if mailerConfig.SkipSSL {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &Mailer{dialer: dialer}
}

// SendEmails via smtp
func (m *Mailer) SendEmails(mailOptions MailOptions) (messageID string, err error) {
	messageID, _ = generateMessageID()

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", mailOptions.EmailFrom, mailOptions.FromName)
	msg.SetHeader("To", mailOptions.EmailTo...)
	msg.SetHeader("Cc", mailOptions.EmailCc...)
	msg.SetHeader("Bcc", mailOptions.EmailBcc...)
	msg.SetHeader("Reply-To", mailOptions.EmailReplyTo...)

	msg.SetHeader("Message-Id", messageID)
	msg.SetHeader("Subject", mailOptions.Subject)
	msg.SetBody("text/html", mailOptions.Content)

	// Send the email to receivers
	if err := m.dialer.DialAndSend(msg); err != nil {
		return "", err
	}
	return messageID, nil
}

var maxBigInt = big.NewInt(math.MaxInt64)

// generateMessageID generates and returns a string suitable for an RFC 2822
// compliant Message-ID, e.g.: <1444789264909237300.3464.1819418242800517193@DESKTOP01>
func generateMessageID() (string, error) {
	t := time.Now().UnixNano()
	pid := os.Getpid()
	rint, err := rand.Int(rand.Reader, maxBigInt)
	if err != nil {
		return "", err
	}
	h, err := os.Hostname()
	// If we can't get the hostname, we'll use localhost
	if err != nil {
		h = "localhost.localdomain"
	}
	msgid := fmt.Sprintf("<%d.%d.%d@%s>", t, pid, rint, h)
	return msgid, nil
}
