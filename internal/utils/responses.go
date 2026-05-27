package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
	}
}

func ValidationError(c *gin.Context, err error) {
	if ve, ok := err.(validator.ValidationErrors); ok && len(ve) > 0 {
		e := ve[0]
		field := strings.ToLower(e.Field())
		var msg string
		switch e.Tag() {
		case "required":
			msg = field + " is required"
		case "email":
			msg = field + " has an invalid format"
		case "min":
			msg = field + " is too short"
		case "max":
			msg = field + " is too long"
		default:
			msg = field + " has an invalid value"
		}
		c.JSON(http.StatusBadRequest, ErrorResponse(msg))
		return
	}

	c.JSON(http.StatusBadRequest, ErrorResponse("invalid request"))
}
