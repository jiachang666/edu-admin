package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("x-request-id")
		if requestID == "" {
			requestID = randomID()
		}
		c.Set("request_id", requestID)
		c.Header("x-request-id", requestID)
		c.Next()
	}
}

func randomID() string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
