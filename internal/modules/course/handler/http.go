package handler

import (
	"edu-admin/internal/app/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func New() *Handler { return &Handler{} }

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.POST("", h.create)
	rg.GET("/options", h.options)
	rg.GET("/:id", h.detail)
	rg.PATCH("/:id", h.update)
}

func (h *Handler) list(c *gin.Context)    { response.Paginated(c, []any{}) }
func (h *Handler) create(c *gin.Context)  { response.Created(c, gin.H{"id": 1}) }
func (h *Handler) options(c *gin.Context) { response.Success(c, []any{}) }
func (h *Handler) detail(c *gin.Context)  { response.Success(c, gin.H{"id": c.Param("id")}) }
func (h *Handler) update(c *gin.Context)  { response.Success(c, gin.H{"updated": true}) }
