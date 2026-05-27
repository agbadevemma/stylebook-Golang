package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Client struct {
	Token   string
	PhoneID string
}

func NewClient(token, phoneId string) *Client {
	return &Client{Token: token, PhoneID: phoneId}
}

type WhatsAppPayload struct {
	MessagingProduct string          `json:"messaging_product"`
	RecipientType    string          `json:"recipient_type"`
	To               string          `json:"to"`
	Type             string          `json:"type"`
	Text             WhatsAppMessage `json:"text"`
}

type WhatsAppMessage struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

func (c *Client) SendText(to, message string) error {

	url := fmt.Sprintf("%v/%v/messages", os.Getenv("FACEBOOK_URL"), c.PhoneID) 
	payload := WhatsAppPayload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "text",
		Text: WhatsAppMessage{
			PreviewURL: false,
			Body:       message,
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	log.Print("An error occured", err)
	if err != nil {

	}
	return err
}
