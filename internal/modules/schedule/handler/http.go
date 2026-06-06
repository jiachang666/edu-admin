package handler

import (
	"edu-admin/internal/app/permission"
	"edu-admin/internal/app/response"
	eduservice "edu-admin/internal/modules/edu/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *eduservice.Service
}

type schedulePayload struct {
	ClassID      uint64 `json:"classId"`
	ScheduleType string `json:"scheduleType"`
	LessonDate   string `json:"lessonDate"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	Classroom    string `json:"classroom"`
	Remark       string `json:"remark"`
}

type scheduleActionPayload struct {
	LessonDate string `json:"lessonDate"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Classroom  string `json:"classroom"`
	Remark     string `json:"remark"`
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
	if !permission.HasFromContext(c, "schedules:view") {
		response.Forbidden(c)
		return
	}

	schedules, scheduleErr := h.service.Schedules()
	if scheduleErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, schedules)
}

func (h *Handler) create(c *gin.Context) {
	if !permission.HasFromContext(c, "schedules:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindSchedulePayload(c)
	if !ok {
		return
	}

	createdItem, createErr := h.service.CreateSchedule(input)
	if createErr != nil {
		response.InternalServerError(c)
		return
	}
	if createdItem.ID == 0 {
		response.Failed(c, 404, "class not found")
		return
	}

	response.Created(c, createdItem)
}

func (h *Handler) detail(c *gin.Context) {
	if !permission.HasFromContext(c, "schedules:view") {
		response.Forbidden(c)
		return
	}

	schedule, found, scheduleErr := h.service.ScheduleDetail(c.Param("id"))
	if scheduleErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "schedule not found")
		return
	}

	response.Success(c, schedule)
}

func (h *Handler) update(c *gin.Context) {
	if !permission.HasFromContext(c, "schedules:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindSchedulePayload(c)
	if !ok {
		return
	}

	updatedItem, found, updateErr := h.service.UpdateSchedule(c.Param("id"), input)
	if updateErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "schedule not found")
		return
	}

	response.Success(c, updatedItem)
}

func (h *Handler) reschedule(c *gin.Context) {
	if !permission.HasFromContext(c, "schedules:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindScheduleActionPayload(c, true)
	if !ok {
		return
	}

	updatedItem, found, updateErr := h.service.Reschedule(c.Param("id"), input)
	if updateErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "schedule not found")
		return
	}

	response.Success(c, updatedItem)
}

func (h *Handler) cancel(c *gin.Context) {
	if !permission.HasFromContext(c, "schedules:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindScheduleActionPayload(c, false)
	if !ok {
		return
	}

	updatedItem, found, updateErr := h.service.CancelSchedule(c.Param("id"), input)
	if updateErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "schedule not found")
		return
	}

	response.Success(c, updatedItem)
}

func (h *Handler) makeup(c *gin.Context) {
	if !permission.HasFromContext(c, "schedules:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindScheduleActionPayload(c, true)
	if !ok {
		return
	}

	createdItem, found, createErr := h.service.CreateMakeupSchedule(c.Param("id"), input)
	if createErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "schedule not found")
		return
	}

	response.Success(c, createdItem)
}

func (h *Handler) attendance(c *gin.Context) {
	if !permission.HasFromContext(c, "attendance:view") {
		response.Forbidden(c)
		return
	}

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
		response.Failed(c, 404, "schedule not found")
		return
	}

	response.Success(c, gin.H{"saved": true})
}

func (h *Handler) homework(c *gin.Context) {
	if !permission.HasFromContext(c, "homeworks:view") {
		response.Forbidden(c)
		return
	}

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
	if !permission.HasFromContext(c, "homeworks:manage") {
		response.Forbidden(c)
		return
	}

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
	if !permission.HasFromContext(c, "homeworks:view") {
		response.Forbidden(c)
		return
	}

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
	if !permission.HasFromContext(c, "homeworks:manage") {
		response.Forbidden(c)
		return
	}

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

func bindSchedulePayload(c *gin.Context) (eduservice.SchedulePayload, bool) {
	var payload schedulePayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "排课表单格式不正确")
		return eduservice.SchedulePayload{}, false
	}

	input := eduservice.SchedulePayload{
		ClassID:      payload.ClassID,
		ScheduleType: strings.TrimSpace(payload.ScheduleType),
		LessonDate:   strings.TrimSpace(payload.LessonDate),
		StartTime:    strings.TrimSpace(payload.StartTime),
		EndTime:      strings.TrimSpace(payload.EndTime),
		Classroom:    strings.TrimSpace(payload.Classroom),
		Remark:       strings.TrimSpace(payload.Remark),
	}

	validationMessage := validateSchedulePayload(input)
	if validationMessage != "" {
		response.Failed(c, 400, validationMessage)
		return eduservice.SchedulePayload{}, false
	}

	return input, true
}

func bindScheduleActionPayload(c *gin.Context, requireTime bool) (eduservice.ScheduleActionPayload, bool) {
	var payload scheduleActionPayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "排课动作表单格式不正确")
		return eduservice.ScheduleActionPayload{}, false
	}

	input := eduservice.ScheduleActionPayload{
		LessonDate: strings.TrimSpace(payload.LessonDate),
		StartTime:  strings.TrimSpace(payload.StartTime),
		EndTime:    strings.TrimSpace(payload.EndTime),
		Classroom:  strings.TrimSpace(payload.Classroom),
		Remark:     strings.TrimSpace(payload.Remark),
	}

	validationMessage := validateScheduleActionPayload(input, requireTime)
	if validationMessage != "" {
		response.Failed(c, 400, validationMessage)
		return eduservice.ScheduleActionPayload{}, false
	}

	return input, true
}

func validateSchedulePayload(input eduservice.SchedulePayload) string {
	if input.ClassID == 0 {
		return "班级不能为空"
	}
	if input.LessonDate == "" {
		return "上课日期不能为空"
	}
	if !isValidDate(input.LessonDate) {
		return "上课日期格式不正确"
	}
	if input.StartTime == "" || input.EndTime == "" {
		return "开始和结束时间不能为空"
	}
	if !isValidClock(input.StartTime) || !isValidClock(input.EndTime) {
		return "开始和结束时间格式不正确"
	}
	if input.Classroom == "" {
		return "教室或地点不能为空"
	}
	if input.StartTime >= input.EndTime {
		return "结束时间必须晚于开始时间"
	}

	return ""
}

func validateScheduleActionPayload(input eduservice.ScheduleActionPayload, requireTime bool) string {
	if requireTime {
		if input.LessonDate == "" {
			return "日期不能为空"
		}
		if !isValidDate(input.LessonDate) {
			return "日期格式不正确"
		}
		if input.StartTime == "" || input.EndTime == "" {
			return "开始和结束时间不能为空"
		}
		if !isValidClock(input.StartTime) || !isValidClock(input.EndTime) {
			return "开始和结束时间格式不正确"
		}
		if input.Classroom == "" {
			return "教室或地点不能为空"
		}
		if input.StartTime >= input.EndTime {
			return "结束时间必须晚于开始时间"
		}
	}

	if !requireTime && input.Remark == "" {
		return "请填写原因说明"
	}

	return ""
}

func isValidDate(value string) bool {
	_, parseErr := time.Parse("2006-01-02", value)
	return parseErr == nil
}

func isValidClock(value string) bool {
	_, parseErr := time.Parse("15:04", value)
	return parseErr == nil
}
