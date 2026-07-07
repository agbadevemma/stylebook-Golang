package otp

import (
	"github.com/gin-gonic/gin"
)

type Module struct {
	handler *Handler
}

func NewModule(otpService *Service) *Module {

	otpHandler := NewHandler(otpService)

	return &Module{handler: otpHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	otp := api.Group("/otp")
	{
		otp.POST("/send", m.handler.SendOtp)
	
	}

}
