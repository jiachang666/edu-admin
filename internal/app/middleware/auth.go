package middleware

import (
	"strings"

	"edu-admin/internal/app/response"
	authservice "edu-admin/internal/modules/auth/service"

	"github.com/gin-gonic/gin"
)

func RequireAuth(authSvc *authservice.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		userID, valid := authSvc.ParseAccessToken(token)
		if !valid {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		profile, found, profileErr := authSvc.AuthProfile(userID)
		if profileErr != nil || !found {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		if profile.Status != "启用" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("current_user_id", profile.ID)
		if len(profile.Roles) > 0 {
			c.Set("current_role", profile.Roles[0])
		}
		c.Set("current_user_name", profile.DisplayName)
		c.Set("current_permissions", profile.Permissions)
		c.Next()
	}
}
