package service

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
	user := api.Group("/services")
	{
		user.POST("/create", m.handler.CreateService)
		user.GET("/list", m.handler.ListServices)
		user.GET("/id", m.handler.GetServiceByID)
		user.PUT("/id", m.handler.UpdateService)
		user.DELETE("/id", m.handler.DeleteService)
	}

}
