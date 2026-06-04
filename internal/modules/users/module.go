package users

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/businesshour"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB, bus business.BusinessWriter, busHour businesshour.BusinessHourWriter) *Module {

	userRepo := NewRepository(db)
	userService := NewService(userRepo, bus, busHour)
	userHadler := NewHandler(userService)

	return &Module{handler: userHadler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	user := api.Group("/auth")
	{
		user.POST("/signup", m.handler.Signup)
		user.POST("/login", m.handler.Login)
	}

}
