package response

import "github.com/gin-gonic/gin"

func Forbidden(c *gin.Context) {
	c.JSON(403, gin.H{
		"code":      403,
		"message":   "forbidden",
		"requestId": c.GetString("request_id"),
	})
}
