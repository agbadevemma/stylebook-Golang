package withdrawal

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/middleware"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/ledger"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB, service *Service, ledgerWriter ledger.LedgerWriter, businesWriter business.BusinessWriter) *Module {
	withdrawalRepo := NewRepository(db)
	withdrawalProvider := NewProvider()
	withdrawalService := NewService(withdrawalRepo, withdrawalProvider, ledgerWriter, businesWriter)
	withdrawalHandler := NewHandler(withdrawalService)

	return &Module{
		handler: withdrawalHandler,
	}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	withdrawal := api.Group("/withdrawal")
	{
		withdrawal.GET("/create", middleware.JWTAuthMiddleware(),
			middleware.RoleMiddleware(utils.Admin), m.handler.CreateWithdrawal)

			
	}
}
