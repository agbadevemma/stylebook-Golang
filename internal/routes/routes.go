package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/whatsapp"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api/v1/")

	whatsappHandler := whatsapp.Module()
	whatappRoutes := api.Group("/whatsapp")
	{
		whatappRoutes.GET("/webhook", whatsappHandler.VerifyWebhook)
		whatappRoutes.POST("/webhook", whatsappHandler.Webhook)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

}
