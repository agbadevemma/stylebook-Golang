package blockedtime

import (
	"time"

	"github.com/google/uuid"
)

type CreateBlockedTimeRequest struct {
	StartTime time.Time `json:"startTime" binding:"required"`
	EndTime   time.Time `json:"endTime" binding:"required"`
	Reason    string    `json:"reason" binding:"required"`
}

type UpdateBlockedTimeRequest struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Reason    string    `json:"reason"`
}

type ResponseBlockedTime struct {
	ID        uuid.UUID `json:"id"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Reason    string    `json:"reason"`
}
