package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/middleware"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Module struct {
	handler *Handler
}

func NewModule(paymentService *Service) *Module {

	paymentHandler := NewHandler(paymentService)

	return &Module{handler: paymentHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	payment := api.Group("/payment")
	{

		payment.POST("/webhook", m.handler.PaymentWebhook)
		payment.GET("/history", middleware.JWTAuthMiddleware(), middleware.RoleMiddleware(utils.Business), m.handler.GetBusinessTransactions)
		payment.GET("/bank/resolve", middleware.JWTAuthMiddleware(), middleware.RoleMiddleware(utils.Business), m.handler.ResolveBankAccount)
		payment.POST(
			"/recipient",
			middleware.JWTAuthMiddleware(),
			middleware.RoleMiddleware(utils.Business),
			m.handler.CreateTransferRecipient,
		)
	}

}
