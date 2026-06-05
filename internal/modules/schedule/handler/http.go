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
	rg.POST("/:id/reschedule", h.reschedule)
	rg.POST("/:id/cancel", h.cancel)
	rg.POST("/:id/makeup", h.makeup)
	rg.GET("/:id/attendance", h.attendance)
	rg.PUT("/:id/attendance", h.saveAttendance)
	rg.GET("/:id/homework", h.homework)
	rg.PUT("/:id/homework", h.saveHomework)
	rg.GET("/:id/feedback", h.feedback)
	rg.PUT("/:id/feedback", h.saveFeedback)
}

func (h *Handler) list(c *gin.Context) {
	response.Paginated(c, demo.Schedules())
}

func (h *Handler) create(c *gin.Context) {
	response.Created(c, gin.H{"id": 4})
}

func (h *Handler) detail(c *gin.Context) {
	schedule, found := demo.FindSchedule(c.Param("id"))
	if !found {
		response.Success(c, gin.H{"id": c.Param("id")})
		return
	}

	response.Success(c, schedule)
}

func (h *Handler) update(c *gin.Context) {
	response.Success(c, gin.H{"updated": true})
}

func (h *Handler) reschedule(c *gin.Context) {
	response.Success(c, gin.H{"rescheduled": true})
}

func (h *Handler) cancel(c *gin.Context) {
	response.Success(c, gin.H{"canceled": true})
}

func (h *Handler) makeup(c *gin.Context) {
	response.Success(c, gin.H{"makeupCreated": true})
}

func (h *Handler) attendance(c *gin.Context) {
	response.Success(c, gin.H{"items": demo.Attendance(c.Param("id"))})
}

func (h *Handler) saveAttendance(c *gin.Context) {
	response.Success(c, gin.H{"saved": true})
}

func (h *Handler) homework(c *gin.Context) {
	response.Success(c, gin.H{"hasHomework": true, "summary": "完成课堂练习册第 12-15 页"})
}

func (h *Handler) saveHomework(c *gin.Context) {
	response.Success(c, gin.H{"saved": true})
}

func (h *Handler) feedback(c *gin.Context) {
	response.Success(c, gin.H{"summary": "课堂状态稳定，建议课后复习错题。"})
}

func (h *Handler) saveFeedback(c *gin.Context) {
	response.Success(c, gin.H{"saved": true})
}
