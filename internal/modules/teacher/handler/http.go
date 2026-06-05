package handler

import (
	"edu-admin/internal/app/response"
	"edu-admin/internal/modules/demo"

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

func (h *Handler) list(c *gin.Context) {
	response.Paginated(c, demo.Teachers())
}

func (h *Handler) create(c *gin.Context) {
	response.Created(c, gin.H{"id": 4})
}

func (h *Handler) options(c *gin.Context) {
	response.Success(c, demo.TeacherOptions())
}

func (h *Handler) detail(c *gin.Context) {
	teacher, found := demo.FindTeacher(c.Param("id"))
	if !found {
		response.Success(c, gin.H{"id": c.Param("id")})
		return
	}

	response.Success(c, teacher)
}

func (h *Handler) update(c *gin.Context) {
	response.Success(c, gin.H{"updated": true})
}
