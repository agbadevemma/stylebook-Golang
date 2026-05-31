package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

func RoleMiddleware(allowedRoles ...utils.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleValue, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusForbidden, utils.ErrorResponse("Role not found"))
			c.Abort()
			return
		}

		userRole := roleValue.(utils.UserRole)

		if slices.Contains(allowedRoles, userRole) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, utils.ErrorResponse("Access Denied"))

		c.Abort()
	}
}
