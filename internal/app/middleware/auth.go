package middleware

import (
	"strings"

	"edu-admin/internal/app/response"

	"github.com/gin-gonic/gin"
)

func RequireAuth(devToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		expected := "Bearer " + devToken
		if strings.TrimSpace(authHeader) != expected {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("current_user_id", uint64(1))
		c.Set("current_role", "super_admin")
		c.Next()
	}
}
