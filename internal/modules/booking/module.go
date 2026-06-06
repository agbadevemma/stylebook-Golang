package booking

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/bookingservice"
	"github.com/stylebookgo/agbadevemma/internal/modules/payment"
	"github.com/stylebookgo/agbadevemma/internal/modules/service"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB, bookingserviceWriter bookingservice.BookingServiceWriter, serviceWriter service.ServiceWriter, paymentWriter payment.PaymentWriter) *Module {

	bookingRepo := NewRepository(db)
	bookingService := NewService(bookingRepo, bookingserviceWriter, serviceWriter, paymentWriter)
	bookingHandler := NewHandler(bookingService)

	return &Module{handler: bookingHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	booking := api.Group("/booking")
	booking.POST("/", m.handler.CreateBooking)
}
