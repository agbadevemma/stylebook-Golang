package users

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/businesshour"
	"github.com/stylebookgo/agbadevemma/internal/modules/otp"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB, bus business.BusinessWriter, busHour businesshour.BusinessHourWriter, otpWriter otp.OtpWriter) *Module {

	userRepo := NewRepository(db)
	userService := NewService(userRepo, bus, busHour, otpWriter)
	userHandler := NewHandler(userService)

	return &Module{handler: userHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	user := api.Group("/auth")
	{
		user.POST("/signup", m.handler.Signup)
		user.POST("/login", m.handler.Login)
		user.POST("/verify", m.handler.VerifyOTP)
		user.POST("/changepassword", m.handler.ChangePassword)
	}

}
