package users

import (
	"github.com/google/uuid"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type SignupRequest struct {
	Name         string         `json:"name" binding:"required"`
	Email        string         `json:"email" binding:"required"`
	Phone        string         `json:"phone" binding:"required"`
	Role         utils.UserRole `json:"role" binding:"required"`
	BusinessName string         `json:"businessName" binding:"required_if=Role admin"` //only used when you want to register a new business with an admin user
	BusinessID   uuid.UUID      `json:"businessId" binding:"required_if=Role barber"`   //only used when you want to register a barber to an existing business
	Password     string         `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
