package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/whatsapp"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	api := r.Group("/api/v1/")

	whatsappHandler := whatsapp.NewModule()
	whatsappHandler.ModuleRoutes(api)

	businessHandler := business.NewModule(db)
	businessHandler.ModuleRoutes(api)

}
