package otp

import "github.com/stylebookgo/agbadevemma/internal/utils"

type SendOTPRequest struct {
	Email string        `json:"email" binding:"required,email"`
	Type  utils.OTPType `json:"type" binding:"required"`
}

type VerifyOTPRequest struct {
	Email string        `json:"email" binding:"required,email"`
	Code  string        `json:"code" binding:"required"`
	Type  utils.OTPType `json:"type" binding:"required"`
}
