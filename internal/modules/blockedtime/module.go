package blockedtime

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

func NewModule(db *gorm.DB, businessWriter business.BusinessWriter) *Module {

	blockedTimeRepo := NewRepository(db)
	blockedTimeService := NewService(blockedTimeRepo, businessWriter)
	blockedTimeHandler := NewHandler(blockedTimeService)

	return &Module{handler: blockedTimeHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	blockedTime := api.Group("/blockedtime")
	{
		blockedTime.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware(utils.Admin))

		blockedTime.POST("/create", m.handler.CreateBlockedTime)
		blockedTime.GET("/list", m.handler.ListBlockedTime)
		blockedTime.GET("/:id", m.handler.GetBlockedTimeByID)
		blockedTime.PUT("/:id", m.handler.UpdateBlockedTime)
		blockedTime.DELETE("/:id", m.handler.DeleteBlockedTime)
	}

}
