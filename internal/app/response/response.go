package response

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
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

func Failed(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"code":      statusCode,
		"message":   message,
		"requestId": c.GetString("request_id"),
	})
}

func InternalServerError(c *gin.Context) {
	Failed(c, 500, "internal server error")
}

func NotImplemented(c *gin.Context) {
	c.JSON(501, gin.H{
		"code":      501,
		"message":   "not implemented yet",
		"requestId": c.GetString("request_id"),
	})
}

func Paginated(c *gin.Context, list any) {
	page := parsePaginationInt(c.Query("page"), defaultPage)
	pageSize := parsePaginationInt(c.Query("pageSize"), defaultPageSize)
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	total, paginatedList := paginateList(list, page, pageSize)

	Success(c, gin.H{
		"list":     paginatedList,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func parsePaginationInt(rawValue string, fallback int) int {
	trimmedValue := strings.TrimSpace(rawValue)
	if trimmedValue == "" {
		return fallback
	}

	parsedValue, parseErr := strconv.Atoi(trimmedValue)
	if parseErr != nil || parsedValue <= 0 {
		return fallback
	}

	return parsedValue
}

func paginateList(list any, page int, pageSize int) (int, any) {
	value := reflect.ValueOf(list)
	if !value.IsValid() {
		return 0, list
	}

	switch value.Kind() {
	case reflect.Array:
		value = value.Slice(0, value.Len())
	case reflect.Slice:
	default:
		return 0, list
	}

	total := value.Len()
	if total == 0 {
		return 0, list
	}

	startIndex := (page - 1) * pageSize
	if startIndex >= total {
		return total, value.Slice(0, 0).Interface()
	}

	endIndex := startIndex + pageSize
	if endIndex > total {
		endIndex = total
	}

	return total, value.Slice(startIndex, endIndex).Interface()
}
