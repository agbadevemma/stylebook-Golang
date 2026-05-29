package service

import "github.com/google/uuid"

type CreateServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
	Duration    int    `json:"duration" binding:"required"`
}

type UpdateServiceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Duration    int    `json:"duration"`
	IsActive    bool   `json:"is_active"`
}

type DeleteServiceRequest struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type ServiceResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Duration    int       `json:"duration"`
	Price       int64     `json:"price"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	BusinessID  uuid.UUID `json:"business_id"`
}
