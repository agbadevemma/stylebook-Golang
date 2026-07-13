package users

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	if req.Role != utils.Business && req.Role != utils.Barber {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid role"))
		return
	}

	err := h.service.Signup(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("User Created", nil))
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}
	token, business, err := h.service.Login(&req)
	
	// Check if error is due to unverified email
	if err != nil && err.Error() == "Email is not verified" {
		c.JSON(http.StatusForbidden, map[string]any{
			"success": false,
			"message": "Email is not verified",
			"data": map[string]string{
				"email": req.Email,
			},
		})
		return
	}
	
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Login Successfully", map[string]any{
		"token":    token,
		"business": business,
	}))

}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	err := h.service.ChangePassword(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(
		"Password Changed successfully",
		nil,
	))
}

func (h *Handler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	err := h.service.VerifyOTP(&req)

	if err != nil {
		log.Printf("Error -->> %v", err.Error())
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("An error  occured"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("OTP verified successfully", nil))
}
