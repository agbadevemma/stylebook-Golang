package users

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB) *Module {

	userRepo := NewRepository(db)
	userService := NewService(userRepo)
	userHadler := NewHandler(userService)

	return &Module{handler: userHadler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	user := api.Group("/auth/user")
	{
		user.GET("/signup", m.handler.Signup)
		user.POST("/login", m.handler.Login)
	}

}
