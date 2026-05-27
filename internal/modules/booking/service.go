package booking

import (
	"github.com/stylebookgo/agbadevemma/internal/modules/client"
)

type Service struct {
	repo         *Repository
	clientWriter client.ClientWriter
}

func NewService(repo *Repository, bus client.ClientWriter) *Service {
	return &Service{repo: repo, clientWriter: bus}
}

func (s *Service) CreateBooking(req *CreateBookingRequest) {

}
