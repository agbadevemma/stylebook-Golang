package whatsapp

import (
	"os"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
}

func NewModule() *Module {
	whatsappFlow := NewFlow()
	whatsappSession := NewSessionService()
	whatsappClient := NewClient(
		os.Getenv("WHATSAPP_TOKEN"),
		os.Getenv("WHATSAPP_PHONE_ID"),
	)
	whatappService := NewService(whatsappFlow, whatsappSession, whatsappClient)
	whatsappHandler := NewHandler(whatappService)
	return &Module{handler: whatsappHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	whatappRoutes := api.Group("/whatsapp")
	{
		whatappRoutes.GET("/webhook", m.handler.VerifyWebhook)
		whatappRoutes.POST("/webhook", m.handler.Webhook)
	}

}
