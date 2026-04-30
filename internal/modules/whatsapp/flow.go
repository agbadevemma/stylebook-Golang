package whatsapp

import (
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Flow struct {
}

func NewFlow() *Flow {
	return &Flow{}
}

func (f *Flow) Handle(session *Session, message string) string {
	switch session.State {

	case string(utils.StateStart):
		session.State = string(utils.StateService)
		return "Welcome 👋 Choose a service: Haircut, Nails, Facial"

	case string(utils.StateService):
		session.Data["service"] = message
		session.State = string(utils.StateStaff)
		return "Choose staff: John, Mary, David"

	case string(utils.StateStaff):
		session.Data["staff"] = message
		session.State = string(utils.StateTime)
		return "Pick a time: 10am, 11am, 2pm"

	case string(utils.StateTime):
		if message == "Yes" {
			session.State = string(utils.StateStart)
			return "Booking confirmed 🎉"
		}
		session.State = string(utils.StateStart)
		return "Booking cancelled"

	default:
		session.State = string(utils.StateStart)
		return "Say hi to start booking 👋"
	}

}
