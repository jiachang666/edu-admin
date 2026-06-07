package handler

import (
	"edu-admin/internal/app/permission"
	"edu-admin/internal/app/response"
	eduservice "edu-admin/internal/modules/edu/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *eduservice.Service
}

func New(service *eduservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
}

func (h *Handler) list(c *gin.Context) {
	if !permission.HasFromContext(c, "operation_logs:view") {
		response.Forbidden(c)
		return
	}

	userID, userErr := parseAuditUintParam(c.Query("userId"))
	if userErr != nil {
		response.Failed(c, 400, "user id is invalid")
		return
	}

	logs, logErr := h.service.OperationLogsWithFilter(eduservice.OperationLogFilter{
		UserID:   userID,
		Module:   strings.TrimSpace(c.Query("module")),
		Action:   strings.TrimSpace(c.Query("action")),
		DateFrom: strings.TrimSpace(c.Query("dateFrom")),
		DateTo:   strings.TrimSpace(c.Query("dateTo")),
	})
	if logErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, logs)
}

func parseAuditUintParam(rawValue string) (uint64, error) {
	trimmedValue := strings.TrimSpace(rawValue)
	if trimmedValue == "" {
		return 0, nil
	}

	return strconv.ParseUint(trimmedValue, 10, 64)
}
