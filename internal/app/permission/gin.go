package permission

import "github.com/gin-gonic/gin"

func HasFromContext(c *gin.Context, target string) bool {
	currentPermissions, exists := c.Get("current_permissions")
	if !exists {
		return false
	}

	permissionList, ok := currentPermissions.([]string)
	if !ok {
		return false
	}

	return Has(permissionList, target)
}
