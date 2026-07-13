package booking

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/middleware"
	"github.com/stylebookgo/agbadevemma/internal/modules/blockedtime"
	"github.com/stylebookgo/agbadevemma/internal/modules/bookingservice"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/businesshour"
	"github.com/stylebookgo/agbadevemma/internal/modules/payment"
	"github.com/stylebookgo/agbadevemma/internal/modules/service"
	"github.com/stylebookgo/agbadevemma/internal/utils"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB, bookingserviceWriter bookingservice.BookingServiceWriter,
	businessHourWriter businesshour.BusinessHourWriter,
	blockedTimesWriter blockedtime.BlockedTimesWriter,
	serviceWriter service.ServiceWriter,
	paymentWriter payment.PaymentWriter, businessWriter business.BusinessWriter) *Module {

	bookingRepo := NewRepository(db)
	bookingService := NewService(bookingRepo, bookingserviceWriter, businessHourWriter, serviceWriter, blockedTimesWriter, paymentWriter, businessWriter)
	bookingHandler := NewHandler(bookingService)

	return &Module{handler: bookingHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	booking := api.Group("/booking")
	booking.POST("/", m.handler.CreateBooking)
	booking.GET("/:bookingId/details", m.handler.GetBookingDetails)

	businessBooking := api.Group("/booking")
	businessBooking.Use(middleware.JWTAuthMiddleware(), middleware.RoleMiddleware(utils.Business))
	businessBooking.GET("/list", m.handler.ListBusinessBookings)
	businessBooking.GET("/stats", m.handler.GetBusinessStats)
	businessBooking.PUT("/:bookingId/reschedule", m.handler.RescheduleBooking)
	businessBooking.PATCH("/:bookingId/status", m.handler.UpdateBookingStatus)
}
