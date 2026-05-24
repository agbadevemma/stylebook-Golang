package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/whatsapp"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	api := r.Group("/api/v1/")

	whatsapp.NewModule().ModuleRoutes(api)

}
