package handler

import (
	"edu-admin/internal/app/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func New() *Handler { return &Handler{} }

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
}

func (h *Handler) list(c *gin.Context) { response.Paginated(c, []any{}) }
