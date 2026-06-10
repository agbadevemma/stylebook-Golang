package ledger

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

	ledgerRepo := NewRepository(db)
	ledgerService := NewService(ledgerRepo)
	ledgerHandler := NewHandler(ledgerService)

	return &Module{handler: ledgerHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	ledger := api.Group("/ledger")
	{
		ledger.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware(utils.Admin))

		ledger.GET("/balance", m.handler.GetBusinessBalance)

	}

}
