package booking

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/bookingservice"
	"github.com/stylebookgo/agbadevemma/internal/modules/service"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB, bookingserviceWriter bookingservice.BookingServiceWriter, serviceWriter service.ServiceWriter) *Module {

	bookingRepo := NewRepository(db)
	bookingService := NewService(bookingRepo, bookingserviceWriter, serviceWriter)
	bookingHandler := NewHandler(bookingService)

	return &Module{handler: bookingHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	booking := api.Group("/booking")
	booking.POST("/", m.handler.CreateBooking)
}
