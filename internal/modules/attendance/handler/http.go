package handler

import (
	"edu-admin/internal/app/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func New() *Handler { return &Handler{} }

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.PATCH("/:id", h.update)
}

func (h *Handler) list(c *gin.Context)   { response.Paginated(c, []any{}) }
func (h *Handler) update(c *gin.Context) { response.Success(c, gin.H{"updated": true}) }
