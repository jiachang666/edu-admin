package handler

import (
	"edu-admin/internal/app/response"
	dashboardservice "edu-admin/internal/modules/dashboard/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *dashboardservice.Service
}

func New(service *dashboardservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/overview", h.overview)
}

func (h *Handler) overview(c *gin.Context) {
	overview, overviewErr := h.service.Overview()
	if overviewErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, overview)
}
