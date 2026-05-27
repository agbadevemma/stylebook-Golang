package users

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB, bus business.BusinessWriter) *Module {

	userRepo := NewRepository(db)
	userService := NewService(userRepo, bus)
	userHadler := NewHandler(userService)

	return &Module{handler: userHadler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	user := api.Group("/auth/user")
	{
		user.POST("/signup", m.handler.Signup)
		user.POST("/login", m.handler.Login)
	}

}
