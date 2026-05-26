package barbers

import (
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(db *gorm.DB) *Module {

	userRepo := NewRepository(db)
	userService := NewService(userRepo)
	userHadler := NewHandler(userService)

	return &Module{handler: userHadler}
}
