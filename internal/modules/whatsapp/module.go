package whatsapp

import (
	"os"
)

func Module() *Handler {
	whatsappFlow := NewFlow()
	whatsappSession := NewSessionService()
	whatsappClient := NewClient(
		os.Getenv("WHATSAPP_TOKEN"),
		os.Getenv("WHATSAPP_PHONE_ID"),
	)
	whatappService := NewService(whatsappFlow, whatsappSession, whatsappClient)
	whatsappHandler := NewHandler(whatappService)
	return whatsappHandler
}
