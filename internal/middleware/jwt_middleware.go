package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/utils"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})
			c.Abort()
			return
		}

		splitToken := strings.Split(authHeader, " ")

		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Access Denied"))
			c.Abort()
			return
		}

		tokenString := splitToken[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid or Expired Token"))
			c.Abort()
			return
		}

		// store user info in context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
