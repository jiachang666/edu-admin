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
	rg.POST("", h.create)
	rg.GET("/:id", h.detail)
	rg.PATCH("/:id", h.update)
	rg.GET("/:id/classes", h.classes)
	rg.POST("/:id/guardians", h.createGuardian)
	rg.PATCH("/:id/guardians/:guardianId", h.updateGuardian)
	rg.DELETE("/:id/guardians/:guardianId", h.deleteGuardian)
}

func (h *Handler) list(c *gin.Context) {
	students, studentErr := h.service.Students()
	if studentErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, students)
}

func (h *Handler) create(c *gin.Context) {
	response.Created(c, gin.H{"id": 5})
}

func (h *Handler) detail(c *gin.Context) {
	student, found, studentErr := h.service.Student(c.Param("id"))
	if studentErr != nil {
		response.InternalServerError(c)
		return
	}
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
	classes, classErr := h.service.StudentClasses(c.Param("id"))
	if classErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, classes)
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
