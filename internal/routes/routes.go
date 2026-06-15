package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/availability"
	"github.com/stylebookgo/agbadevemma/internal/modules/blockedtime"
	"github.com/stylebookgo/agbadevemma/internal/modules/booking"
	"github.com/stylebookgo/agbadevemma/internal/modules/bookingservice"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/businesshour"
	"github.com/stylebookgo/agbadevemma/internal/modules/ledger"
	"github.com/stylebookgo/agbadevemma/internal/modules/notification/email"
	"github.com/stylebookgo/agbadevemma/internal/modules/otp"
	"github.com/stylebookgo/agbadevemma/internal/modules/payment"
	"github.com/stylebookgo/agbadevemma/internal/modules/service"
	"github.com/stylebookgo/agbadevemma/internal/modules/users"
	"github.com/stylebookgo/agbadevemma/internal/modules/withdrawal"
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
	ledgerRepo := ledger.NewRepository(db)
	paymentProvider := payment.NewProvider()
	withdrawalRepo := withdrawal.NewRepository(db)
	withdrawalProvider := withdrawal.NewProvider()
	withdrawalService := withdrawal.NewService(withdrawalRepo, withdrawalProvider, ledgerRepo, businessRepo)
	emailService := email.NewEmailService()
	otpRepo := otp.NewRepository(db)
	otpService := otp.NewService(otpRepo, emailService)
	otp.NewModule(otpService).ModuleRoutes(api)
	paymentService := payment.NewService(paymentRepo, paymentProvider, ledgerRepo, businessRepo, withdrawalService)
	bookingServiceRepo := bookingservice.NewRepository(db)

	users.NewModule(db, businessRepo, businessHourRepo, otpService).ModuleRoutes(api)
	service.NewModule(db, businessRepo).ModuleRoutes(api)
	payment.NewModule(paymentService).ModuleRoutes(api)
	withdrawal.NewModule(db, ledgerRepo, businessRepo).ModuleRoutes(api)
	businesshour.NewModule(db).ModuleRoutes(api)
	blockedtime.NewModule(db, businessRepo).ModuleRoutes(api)
	booking.NewModule(db, bookingServiceRepo, businessHourRepo, blockedRepo, serviceRepo, paymentService, businessRepo).ModuleRoutes(api)
	availability.NewModule(bookingRepo, blockedRepo, businessHourRepo, businessRepo).ModuleRoutes(api)
	ledger.NewModule(db).ModuleRoutes(api)
}
