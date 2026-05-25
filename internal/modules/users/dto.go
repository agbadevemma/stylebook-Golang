package users

import "github.com/stylebookgo/agbadevemma/internal/utils"

type SignupRequest struct {
	Name         string         `json:"name" binding:"required"`
	Email        string         `json:"email" binding:"required"`
	Phone        string         `json:"phone" binding:"required"`
	Role         utils.UserRole `json:"role" binding:"required"`
	BusinessName string         `json:"businessName" binding:"required"`
	Password     string         `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
