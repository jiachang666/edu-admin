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
	h.registerHomeworkRoutes(rg)
}

func (h *Handler) RegisterFeedbackRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.feedbacks)
}

func (h *Handler) registerHomeworkRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.GET("/feedbacks", h.feedbacks)
}

func (h *Handler) list(c *gin.Context) {
	if !permission.HasFromContext(c, "homeworks:view") {
		response.Forbidden(c)
		return
	}

	scope, scopeErr := h.service.ScopeForUser(c.GetUint64("current_user_id"), c.GetString("current_role"))
	if scopeErr != nil {
		response.InternalServerError(c)
		return
	}

	classID, classErr := parseHomeworkUintParam(c.Query("classId"))
	if classErr != nil {
		response.Failed(c, 400, "class id is invalid")
		return
	}

	teacherID, teacherErr := parseHomeworkUintParam(c.Query("teacherId"))
	if teacherErr != nil {
		response.Failed(c, 400, "teacher id is invalid")
		return
	}

	items, itemErr := h.service.HomeworksWithFilter(eduservice.HomeworkFilter{
		ClassID:   classID,
		TeacherID: teacherID,
		DateFrom:  strings.TrimSpace(c.Query("dateFrom")),
		DateTo:    strings.TrimSpace(c.Query("dateTo")),
		Scope:     scope,
	})
	if itemErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, items)
}

func (h *Handler) feedbacks(c *gin.Context) {
	if !permission.HasFromContext(c, "homeworks:view") {
		response.Forbidden(c)
		return
	}

	scope, scopeErr := h.service.ScopeForUser(c.GetUint64("current_user_id"), c.GetString("current_role"))
	if scopeErr != nil {
		response.InternalServerError(c)
		return
	}

	classID, classErr := parseHomeworkUintParam(c.Query("classId"))
	if classErr != nil {
		response.Failed(c, 400, "class id is invalid")
		return
	}

	teacherID, teacherErr := parseHomeworkUintParam(c.Query("teacherId"))
	if teacherErr != nil {
		response.Failed(c, 400, "teacher id is invalid")
		return
	}

	items, itemErr := h.service.FeedbacksWithFilter(eduservice.FeedbackFilter{
		ClassID:   classID,
		TeacherID: teacherID,
		DateFrom:  strings.TrimSpace(c.Query("dateFrom")),
		DateTo:    strings.TrimSpace(c.Query("dateTo")),
		Scope:     scope,
	})
	if itemErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, items)
}

func parseHomeworkUintParam(rawValue string) (uint64, error) {
	trimmedValue := strings.TrimSpace(rawValue)
	if trimmedValue == "" {
		return 0, nil
	}

	return strconv.ParseUint(trimmedValue, 10, 64)
}
