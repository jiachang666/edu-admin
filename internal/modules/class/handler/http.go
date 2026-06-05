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
	rg.GET("/:id/students", h.students)
	rg.POST("/:id/students", h.addStudents)
	rg.DELETE("/:id/students/:studentId", h.removeStudent)
	rg.GET("/:id/schedules/upcoming", h.upcomingSchedules)
}

func (h *Handler) list(c *gin.Context) {
	response.Paginated(c, demo.Classes())
}

func (h *Handler) create(c *gin.Context) {
	response.Created(c, gin.H{"id": 4})
}

func (h *Handler) detail(c *gin.Context) {
	classItem, found := demo.FindClass(c.Param("id"))
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
	response.Success(c, demo.ClassStudents(c.Param("id")))
}

func (h *Handler) addStudents(c *gin.Context) {
	response.Success(c, gin.H{"added": true})
}

func (h *Handler) removeStudent(c *gin.Context) {
	response.Success(c, gin.H{"removed": true})
}

func (h *Handler) upcomingSchedules(c *gin.Context) {
	response.Success(c, demo.UpcomingSchedules(c.Param("id")))
}
