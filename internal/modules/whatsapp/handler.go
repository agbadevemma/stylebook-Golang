package whatsapp

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type WhatsAppWebhookPayload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Value Value `json:"value"`
}

type Value struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	From string `json:"from"`
	Text Text   `json:"text"`
	Type string `json:"type"`
}

type Text struct {
	Body string `json:"body"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}

}

func (h *Handler) Webhook(c *gin.Context) {

	var payload WhatsAppWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Print("err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	fmt.Println("payload", payload)
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {

			if len(change.Value.Messages) == 0 {
				continue
			}

			msg := change.Value.Messages[0]
			log.Print("msg.Text.Body", msg.From)
			h.service.HandleIncomingMessage(msg.From, msg.Text.Body)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

func (h *Handler) VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub_mode")
	token := c.Query("hub_verify_token")
	challenge := c.Query("hub_challenge")

	if mode == "subscribe" && token == os.Getenv("WHATSAPP_VERIFY_TOKEN") {
		c.String(200, challenge)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "verification failed"})
}
