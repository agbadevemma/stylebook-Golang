package payment

import (
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
}

func NewModule(paymentService *Service) *Module {

	paymentHandler := NewHandler(paymentService)

	return &Module{handler: paymentHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	service := api.Group("/payment")
	{

		service.POST("/webhook", m.handler.PaystackWebhook)

	}

}
