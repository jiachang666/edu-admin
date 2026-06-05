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
	if !permission.HasFromContext(c, "homeworks:view") {
		response.Forbidden(c)
		return
	}

	items, itemErr := h.service.Homeworks()
	if itemErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, items)
}
