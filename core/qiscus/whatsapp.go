package qiscus

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type QiscusConfig struct {
	Namespace string
	AppID     string
	SecretKey string
	Toggle    bool
}
type Qiscus struct {
	namespace string
	appID     string
	secretKey string
	toggle    bool
}

func NewQiscus(config QiscusConfig) *Qiscus {
	return &Qiscus{namespace: config.Namespace, appID: config.AppID, secretKey: config.SecretKey, toggle: config.Toggle}
}

func (q *Qiscus) GetNamespace() string {
	return q.namespace
}

func (q *Qiscus) SendWhatsapp(content string) error {
	if !q.toggle {
		return nil
	}
	waURL := fmt.Sprintf("https://multichannel.qiscus.com/whatsapp/v1/%s/4022/messages", q.appID)
	req, err := http.NewRequest("POST", waURL, bytes.NewBuffer([]byte(content)))
	if err != nil {
		return fmt.Errorf("Failed SendWhatsapp because %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Qiscus-App-Id", q.appID)
	req.Header.Set("Qiscus-Secret-Key", q.secretKey)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed SendWhatsapp because %s", err)
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Failed SendWhatsapp because %s", err)
	}
	bodyString := string(bodyBytes)
	if res.StatusCode == http.StatusOK {
		fmt.Println("SendWhatsapp response :", bodyString)
	} else {
		return fmt.Errorf("Failed SendWhatsapp because %s", bodyString)
	}
	return nil
}
