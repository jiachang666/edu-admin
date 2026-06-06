package handler

import (
	"edu-admin/internal/app/permission"
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
	if !permission.HasFromContext(c, "dashboard:view") {
		response.Forbidden(c)
		return
	}

	scope, scopeErr := h.service.ScopeForUser(c.GetUint64("current_user_id"), c.GetString("current_role"))
	if scopeErr != nil {
		response.InternalServerError(c)
		return
	}

	overview, overviewErr := h.service.OverviewWithScope(scope)
	if overviewErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, overview)
}
