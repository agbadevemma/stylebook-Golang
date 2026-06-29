package business

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/middleware"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{
		handler: handler,
	}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	business := api.Group("/business")
	{
		business.GET(
			"/details",
			middleware.JWTAuthMiddleware(),
			middleware.RoleMiddleware(utils.Business),
			m.handler.GetBusinessSettings,
		)
		business.PUT(
			"/update",
			middleware.JWTAuthMiddleware(),
			middleware.RoleMiddleware(utils.Business),
			m.handler.UpdateBusinessDetails,
		)
		business.POST(
			"/withdrawal-pin/setup",
			middleware.JWTAuthMiddleware(),
			middleware.RoleMiddleware(utils.Business),
			m.handler.SetupWithdrawalPIN,
		)
		business.PUT(
			"/withdrawal-pin/update",
			middleware.JWTAuthMiddleware(),
			middleware.RoleMiddleware(utils.Business),
			m.handler.UpdateWithdrawalPIN,
		)
	}
}
