package blockedtime

import "time"

type CreateBlockedTimeRequest struct {
	StartTime time.Time `json:"startTime" binding:"required"`
	EndTime   time.Time `json:"endTime" binding:"required"`
	Reason    string    `json:"reason"`
}

type UpdateBlockedTimeRequest struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Reason    string    `json:"reason"`
}
