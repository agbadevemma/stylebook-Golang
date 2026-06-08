package service

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/middleware"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB, bus business.BusinessWriter) *Module {

	serviceRepo := NewRepository(db)
	serviceService := NewService(serviceRepo, bus)
	serviceHandler := NewHandler(serviceService)

	return &Module{handler: serviceHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	service := api.Group("/services")
	{
		service.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware(utils.Admin))

		service.POST("/create", m.handler.CreateService)
		service.GET("/list", m.handler.ListServices)
		service.GET("/:id", m.handler.GetServiceByID)
		service.PUT("/:id", m.handler.UpdateService)
		service.DELETE("/:id", m.handler.DeleteService)
	}

}
