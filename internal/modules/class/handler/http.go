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
	rg.GET("/:id/students", h.students)
	rg.POST("/:id/students", h.addStudents)
	rg.DELETE("/:id/students/:studentId", h.removeStudent)
	rg.GET("/:id/schedules/upcoming", h.upcomingSchedules)
}

func (h *Handler) list(c *gin.Context) {
	classes, classErr := h.service.Classes()
	if classErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, classes)
}

func (h *Handler) create(c *gin.Context) {
	response.Created(c, gin.H{"id": 4})
}

func (h *Handler) detail(c *gin.Context) {
	classItem, found, classErr := h.service.Class(c.Param("id"))
	if classErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Success(c, gin.H{"id": c.Param("id")})
		return
	}

	response.Success(c, classItem)
}

func (h *Handler) update(c *gin.Context) {
	response.Success(c, gin.H{"updated": true})
}

func (h *Handler) students(c *gin.Context) {
	students, studentErr := h.service.ClassStudents(c.Param("id"))
	if studentErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, students)
}

func (h *Handler) addStudents(c *gin.Context) {
	response.Success(c, gin.H{"added": true})
}

func (h *Handler) removeStudent(c *gin.Context) {
	response.Success(c, gin.H{"removed": true})
}

func (h *Handler) upcomingSchedules(c *gin.Context) {
	schedules, scheduleErr := h.service.UpcomingSchedules(c.Param("id"))
	if scheduleErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, schedules)
}
