package business

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB) *Module {
	return &Module{handler: NewHandler(NewService(NewRepo(db)))}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	businessRoutes := api.Group("/business")
	{
		businessRoutes.GET("/", m.handler.GetBusiness)
	}
}
