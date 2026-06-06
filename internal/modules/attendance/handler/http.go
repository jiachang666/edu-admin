package handler

import (
	"edu-admin/internal/app/permission"
	"edu-admin/internal/app/response"
	eduservice "edu-admin/internal/modules/edu/service"
	"strconv"
	"strings"

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
	rg.PATCH("/:id", h.update)
}

func (h *Handler) list(c *gin.Context) {
	if !permission.HasFromContext(c, "attendance:view") {
		response.Forbidden(c)
		return
	}

	mode := strings.TrimSpace(c.Query("mode"))
	classID, classErr := parseUintParam(c.Query("classId"))
	if classErr != nil {
		response.Failed(c, 400, "class id is invalid")
		return
	}

	studentID, studentErr := parseUintParam(c.Query("studentId"))
	if studentErr != nil {
		response.Failed(c, 400, "student id is invalid")
		return
	}

	dateFilter := c.Query("date")
	statusFilter := c.Query("status")
	if mode == "records" || classID > 0 || studentID > 0 || dateFilter != "" || statusFilter != "" {
		items, itemErr := h.service.AttendanceRecords(eduservice.AttendanceRecordFilter{
			ClassID:   classID,
			StudentID: studentID,
			Date:      dateFilter,
			Status:    statusFilter,
		})
		if itemErr != nil {
			response.InternalServerError(c)
			return
		}

		response.Paginated(c, items)
		return
	}

	items, itemErr := h.service.AttendanceSessions()
	if itemErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, items)
}

func (h *Handler) update(c *gin.Context) {
	if !permission.HasFromContext(c, "attendance:manage") {
		response.Forbidden(c)
		return
	}

	var payload eduservice.AttendanceSavePayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "invalid attendance payload")
		return
	}

	saved, saveErr := h.service.SaveAttendance(c.Param("id"), payload, c.GetString("current_user_name"))
	if saveErr != nil {
		response.InternalServerError(c)
		return
	}
	if !saved {
		response.Failed(c, 404, "attendance session not found")
		return
	}

	response.Success(c, gin.H{"saved": true})
}

func parseUintParam(rawValue string) (uint64, error) {
	trimmedValue := strings.TrimSpace(rawValue)
	if trimmedValue == "" {
		return 0, nil
	}

	return strconv.ParseUint(trimmedValue, 10, 64)
}
