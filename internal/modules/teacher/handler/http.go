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
	rg.GET("/options", h.options)
	rg.GET("/:id", h.detail)
	rg.PATCH("/:id", h.update)
}

func (h *Handler) list(c *gin.Context) {
	teachers, teacherErr := h.service.Teachers()
	if teacherErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, teachers)
}

func (h *Handler) create(c *gin.Context) {
	response.Created(c, gin.H{"id": 4})
}

func (h *Handler) options(c *gin.Context) {
	options, optionErr := h.service.TeacherOptions()
	if optionErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, options)
}

func (h *Handler) detail(c *gin.Context) {
	teacher, found, teacherErr := h.service.Teacher(c.Param("id"))
	if teacherErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Success(c, gin.H{"id": c.Param("id")})
		return
	}

	response.Success(c, teacher)
}

func (h *Handler) update(c *gin.Context) {
	response.Success(c, gin.H{"updated": true})
}
