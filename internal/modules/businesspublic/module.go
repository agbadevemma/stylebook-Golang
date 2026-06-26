package businesspublic

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
)

type Module struct {
	handler *Handler
}

func NewModule(businessWriter business.BusinessWriter) *Module {
	businessService := NewService(businessWriter)
	businessHandler := NewHandler(businessService)
	return &Module{handler: businessHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	businessRoute := api.Group("/business")
	{
		businessRoute.GET("/public/:id", m.handler.GetPublicBusiness)
	}
}
