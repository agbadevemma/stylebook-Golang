package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/users"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	api := r.Group("/api/v1/")
	businessRepo := business.NewRepository(db)
	users.NewModule(db, businessRepo).ModuleRoutes(api)

}
