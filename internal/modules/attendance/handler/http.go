package handler

import (
	"edu-admin/internal/app/permission"
	"edu-admin/internal/app/response"
	eduservice "edu-admin/internal/modules/edu/service"
	"errors"
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

	scope, scopeErr := h.service.ScopeForUser(c.GetUint64("current_user_id"), c.GetString("current_role"))
	if scopeErr != nil {
		response.InternalServerError(c)
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
	dateFromFilter := strings.TrimSpace(c.Query("dateFrom"))
	dateToFilter := strings.TrimSpace(c.Query("dateTo"))
	statusFilter := c.Query("status")
	if mode == "records" || classID > 0 || studentID > 0 || dateFilter != "" || dateFromFilter != "" || dateToFilter != "" || statusFilter != "" {
		items, itemErr := h.service.AttendanceRecords(eduservice.AttendanceRecordFilter{
			ClassID:   classID,
			StudentID: studentID,
			Date:      dateFilter,
			DateFrom:  dateFromFilter,
			DateTo:    dateToFilter,
			Status:    statusFilter,
			Scope:     scope,
		})
		if itemErr != nil {
			response.InternalServerError(c)
			return
		}

		response.Paginated(c, items)
		return
	}

	items, itemErr := h.service.AttendanceSessionsWithScope(scope)
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

	scope, scopeErr := h.service.ScopeForUser(c.GetUint64("current_user_id"), c.GetString("current_role"))
	if scopeErr != nil {
		response.InternalServerError(c)
		return
	}

	recordItem, found, recordErr := h.service.AttendanceRecord(c.Param("id"))
	if recordErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found || !h.service.ScopeAllowsTeacher(scope, recordItem.TeacherID) {
		response.Failed(c, 404, "attendance record not found")
		return
	}

	var payload eduservice.AttendanceRecordUpdatePayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "invalid attendance payload")
		return
	}

	updatedItem, updateFound, updateErr := h.service.UpdateAttendanceRecord(c.Param("id"), payload, currentOperator(c))
	if errors.Is(updateErr, eduservice.ErrInvalidAttendanceStatus) {
		response.Failed(c, 400, "attendance status is invalid")
		return
	}
	if updateErr != nil {
		response.InternalServerError(c)
		return
	}
	if !updateFound {
		response.Failed(c, 404, "attendance record not found")
		return
	}

	response.Success(c, updatedItem)
}

func parseUintParam(rawValue string) (uint64, error) {
	trimmedValue := strings.TrimSpace(rawValue)
	if trimmedValue == "" {
		return 0, nil
	}

	return strconv.ParseUint(trimmedValue, 10, 64)
}

func currentOperator(c *gin.Context) eduservice.Operator {
	return eduservice.Operator{
		UserID:      c.GetUint64("current_user_id"),
		DisplayName: c.GetString("current_user_name"),
	}
}
