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
	schedules, scheduleErr := h.service.Schedules()
	if scheduleErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, schedules)
}

func (h *Handler) create(c *gin.Context) {
	response.Created(c, gin.H{"id": 4})
}

func (h *Handler) detail(c *gin.Context) {
	schedule, found, scheduleErr := h.service.Schedule(c.Param("id"))
	if scheduleErr != nil {
		response.InternalServerError(c)
		return
	}
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
	scheduleItem, found, scheduleErr := h.service.Schedule(c.Param("id"))
	if scheduleErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "schedule not found")
		return
	}

	items, itemErr := h.service.Attendance(c.Param("id"))
	if itemErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, gin.H{"schedule": scheduleItem, "items": items})
}

func (h *Handler) saveAttendance(c *gin.Context) {
	var payload eduservice.AttendanceSavePayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "invalid attendance payload")
		return
	}

	saved, saveErr := h.service.SaveAttendance(c.Param("id"), payload)
	if saveErr != nil {
		response.InternalServerError(c)
		return
	}
	if !saved {
		response.Failed(c, 404, "schedule not found")
		return
	}

	response.Success(c, gin.H{"saved": true})
}

func (h *Handler) homework(c *gin.Context) {
	item, found, itemErr := h.service.Homework(c.Param("id"))
	if itemErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Success(c, gin.H{"scheduleId": c.Param("id")})
		return
	}

	response.Success(c, item)
}

func (h *Handler) saveHomework(c *gin.Context) {
	var payload eduservice.HomeworkPayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "invalid homework payload")
		return
	}

	item, found, saveErr := h.service.SaveHomework(c.Param("id"), payload)
	if saveErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "schedule not found")
		return
	}

	response.Success(c, item)
}

func (h *Handler) feedback(c *gin.Context) {
	item, found, itemErr := h.service.Feedback(c.Param("id"))
	if itemErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Success(c, gin.H{"scheduleId": c.Param("id")})
		return
	}

	response.Success(c, item)
}

func (h *Handler) saveFeedback(c *gin.Context) {
	var payload eduservice.FeedbackPayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "invalid feedback payload")
		return
	}

	item, found, saveErr := h.service.SaveFeedback(c.Param("id"), payload)
	if saveErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "schedule not found")
		return
	}

	response.Success(c, item)
}
