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
	rg.GET("/:id", h.detail)
	rg.PATCH("/:id", h.update)
	rg.GET("/:id/classes", h.classes)
	rg.POST("/:id/guardians", h.createGuardian)
	rg.PATCH("/:id/guardians/:guardianId", h.updateGuardian)
	rg.DELETE("/:id/guardians/:guardianId", h.deleteGuardian)
}

func (h *Handler) list(c *gin.Context) {
	response.Paginated(c, demo.Students())
}

func (h *Handler) create(c *gin.Context) {
	response.Created(c, gin.H{"id": 5})
}

func (h *Handler) detail(c *gin.Context) {
	student, found := demo.FindStudent(c.Param("id"))
	if !found {
		response.Success(c, gin.H{"id": c.Param("id")})
		return
	}

	response.Success(c, student)
}

func (h *Handler) update(c *gin.Context) {
	response.Success(c, gin.H{"updated": true})
}

func (h *Handler) classes(c *gin.Context) {
	response.Success(c, demo.StudentClasses(c.Param("id")))
}

func (h *Handler) createGuardian(c *gin.Context) {
	response.Created(c, gin.H{"id": 1})
}

func (h *Handler) updateGuardian(c *gin.Context) {
	response.Success(c, gin.H{"updated": true})
}

func (h *Handler) deleteGuardian(c *gin.Context) {
	response.Success(c, gin.H{"deleted": true})
}
