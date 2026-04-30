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

func (c *Client) SendText(to, message string) error {
	url := fmt.Sprintf("%v/%v/messages", os.Getenv("FACEBOOK_URL"), c.PhoneID)
	fmt.Println("url", url)
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "text",
		"text": map[string]any{
			"preview_url": false,
			"body":        message,
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
