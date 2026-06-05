package handler

import (
	"edu-admin/internal/app/permission"
	"edu-admin/internal/app/response"
	eduservice "edu-admin/internal/modules/edu/service"

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
	currentPermissions, exists := c.Get("current_permissions")
	if !exists {
		response.Forbidden(c)
		return
	}

	permissionList, ok := currentPermissions.([]string)
	if !ok || !permission.Has(permissionList, "operation_logs:view") {
		response.Forbidden(c)
		return
	}

	logs, logErr := h.service.OperationLogs()
	if logErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, logs)
}
