package response

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"code":      0,
		"message":   "ok",
		"data":      data,
		"requestId": c.GetString("request_id"),
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(201, gin.H{
		"code":      0,
		"message":   "ok",
		"data":      data,
		"requestId": c.GetString("request_id"),
	})
}

func Unauthorized(c *gin.Context) {
	c.JSON(401, gin.H{
		"code":      401,
		"message":   "unauthorized",
		"requestId": c.GetString("request_id"),
	})
}

func NotImplemented(c *gin.Context) {
	c.JSON(501, gin.H{
		"code":      501,
		"message":   "not implemented yet",
		"requestId": c.GetString("request_id"),
	})
}

func Paginated(c *gin.Context, list any) {
	total := 0
	value := reflect.ValueOf(list)
	if value.IsValid() {
		switch value.Kind() {
		case reflect.Array, reflect.Slice:
			total = value.Len()
		}
	}

	Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     1,
		"pageSize": 20,
	})
}
