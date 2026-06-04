package businesshour

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

	businessHourRepo := NewRepository(db)
	businessHourService := NewService(businessHourRepo)
	businessHourHandler := NewHandler(businessHourService)

	return &Module{handler: businessHourHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	businessHour := api.Group("/business")
	{
		businessHour.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware(utils.Admin))
		businessHour.GET("/hours", m.handler.ListBusinessHours)
		businessHour.POST("/hours", m.handler.CreateBusinessHour)
		businessHour.PUT("/hours", m.handler.UpdateBusinessHour)
	}
}
