package handler

import (
	"edu-admin/internal/app/response"
	eduservice "edu-admin/internal/modules/edu/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *eduservice.Service
}

type classStudentPayload struct {
	StudentIDs []uint64 `json:"studentIds"`
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
	classDetail, found, classErr := h.service.ClassDetail(c.Param("id"))
	if classErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "class not found")
		return
	}

	response.Success(c, classDetail)
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
	var payload classStudentPayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "班级学员表单格式不正确")
		return
	}

	studentIDs := make([]uint64, 0, len(payload.StudentIDs))
	seenIDs := make(map[uint64]bool, len(payload.StudentIDs))
	for _, studentID := range payload.StudentIDs {
		if studentID == 0 || seenIDs[studentID] {
			continue
		}

		studentIDs = append(studentIDs, studentID)
		seenIDs[studentID] = true
	}

	if len(studentIDs) == 0 {
		response.Failed(c, 400, "至少选择一个学员")
		return
	}

	added, addErr := h.service.AddStudentsToClass(c.Param("id"), studentIDs)
	if addErr != nil {
		response.InternalServerError(c)
		return
	}
	if !added {
		response.Failed(c, 404, "class not found")
		return
	}

	response.Success(c, gin.H{"added": true})
}

func (h *Handler) removeStudent(c *gin.Context) {
	studentID := strings.TrimSpace(c.Param("studentId"))
	if studentID == "" {
		response.Failed(c, 400, "student id is required")
		return
	}

	removed, removeErr := h.service.RemoveStudentFromClass(c.Param("id"), studentID)
	if removeErr != nil {
		response.InternalServerError(c)
		return
	}
	if !removed {
		response.Failed(c, 404, "class student not found")
		return
	}

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
