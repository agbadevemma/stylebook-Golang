package availability

import (
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/modules/blockedtime"
	"github.com/stylebookgo/agbadevemma/internal/modules/booking"
	"github.com/stylebookgo/agbadevemma/internal/modules/business"
	"github.com/stylebookgo/agbadevemma/internal/modules/businesshour"
)

type Module struct {
	handler *Handler
}

func NewModule(bookingWriter booking.BookingWriter, blockedTimesWriter blockedtime.BlockedTimesWriter, businessHourWriter businesshour.BusinessHourWriter, businessWriter business.BusinessWriter) *Module {

	availabilityService := NewService(bookingWriter, blockedTimesWriter, businessHourWriter, businessWriter)
	availabilityHandler := NewHandler(availabilityService)

	return &Module{handler: availabilityHandler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	availability := api.Group("/availability")
	{

		availability.GET("/range", m.handler.GetAvailability)
	}
}
