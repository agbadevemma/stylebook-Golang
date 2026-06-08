package businesshour

import (
	"time"

	"github.com/google/uuid"
)

type CreateBusinessHourRequest struct {
	DayOfWeek int    `json:"dayOfWeek" binding:"required,min=0,max=6"`
	IsOpen    bool   `json:"isOpen" binding:"required"`
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

type UpdateBusinessHourRequest struct {
	DayOfWeek *int   `json:"dayOfWeek" binding:"required,min=0,max=6"`
	IsOpen    bool   `json:"isOpen"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type UpdateBusinessHourRequestWrapper struct {
	Hours []UpdateBusinessHourRequest `json:"hours" binding:"required,dive"`
}
type BusinessHoursResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DayOfWeek int       `json:"dayOfWeek"`
	IsOpen    bool      `json:"isOpen"`
	StartTime string    `json:"startTime"`
	EndTime   string    `json:"endTime"`
}
