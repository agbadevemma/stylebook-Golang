package businesshour

type CreateBusinessHourRequest struct {
	DayOfWeek int    `json:"dayOfWeek" binding:"required,min=0,max=6"`
	IsOpen    bool   `json:"isOpen" binding:"required"`
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

type UpdateBusinessHourRequest struct {
	DayOfWeek int    `json:"dayOfWeek" binding:"required,min=0,max=6"`
	IsOpen    bool   `json:"isOpen"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
