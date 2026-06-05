package handler

import (
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
	rg.PATCH("/:id", h.update)
}

func (h *Handler) list(c *gin.Context) {
	items, itemErr := h.service.AttendanceSessions()
	if itemErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, items)
}

func (h *Handler) update(c *gin.Context) {
	var payload eduservice.AttendanceSavePayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "invalid attendance payload")
		return
	}

	saved, saveErr := h.service.SaveAttendance(c.Param("id"), payload)
	if saveErr != nil {
		response.InternalServerError(c)
		return
	}
	if !saved {
		response.Failed(c, 404, "attendance session not found")
		return
	}

	response.Success(c, gin.H{"saved": true})
}
