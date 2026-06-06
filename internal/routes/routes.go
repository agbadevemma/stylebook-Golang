package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/availability"
	"github.com/stylebookgo/agbadevemma/internal/modules/blockedtime"
	"github.com/stylebookgo/agbadevemma/internal/modules/booking"
	"github.com/stylebookgo/agbadevemma/internal/modules/bookingservice"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/businesshour"
	"github.com/stylebookgo/agbadevemma/internal/modules/payment"
	"github.com/stylebookgo/agbadevemma/internal/modules/service"
	"github.com/stylebookgo/agbadevemma/internal/modules/users"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	api := r.Group("/api/v1/")
	businessRepo := business.NewRepository(db)
	businessHourRepo := businesshour.NewRepository(db)
	blockedRepo := blockedtime.NewRepository(db)
	bookingRepo := booking.NewRepository(db)
	serviceRepo := service.NewRepository(db)
	paymentRepo := payment.NewRepository(db)
	paymentProvider := payment.NewProvider()
	paymentService := payment.NewService(paymentRepo, paymentProvider)
	bookingServiceRepo := bookingservice.NewRepository(db)

	users.NewModule(db, businessRepo, businessHourRepo).ModuleRoutes(api)
	service.NewModule(db).ModuleRoutes(api)
	payment.NewModule(paymentService).ModuleRoutes(api)
	businesshour.NewModule(db).ModuleRoutes(api)
	blockedtime.NewModule(db).ModuleRoutes(api)
	booking.NewModule(db, bookingServiceRepo, serviceRepo, paymentService).ModuleRoutes(api)
	availability.NewModule(bookingRepo, blockedRepo, businessHourRepo, businessRepo).ModuleRoutes(api)

}
