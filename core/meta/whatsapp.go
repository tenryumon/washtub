package meta

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type WhatsappConfig struct {
	Token         string
	PhoneNumberID string
}
type Whatsapp struct {
	token         string
	phoneNumberID string
}

func NewWhatsapp(config WhatsappConfig) *Whatsapp {
	return &Whatsapp{token: config.Token, phoneNumberID: config.PhoneNumberID}
}

func (wa *Whatsapp) SendWhatsapp(content string) error {
	waURL := fmt.Sprintf("https://graph.facebook.com/v15.0/%s/messages", wa.phoneNumberID)
	req, err := http.NewRequest("POST", waURL, bytes.NewBuffer([]byte(content)))
	if err != nil {
		return fmt.Errorf("Failed SendWhatsapp because %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", wa.token)
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
