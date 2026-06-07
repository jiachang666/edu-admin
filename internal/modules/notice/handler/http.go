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

type noticePayload struct {
	Title             string   `json:"title"`
	Content           string   `json:"content"`
	Category          string   `json:"category"`
	TargetScope       string   `json:"targetScope"`
	RelatedClassID    uint64   `json:"relatedClassId"`
	RelatedScheduleID uint64   `json:"relatedScheduleId"`
	StudentIDs        []uint64 `json:"studentIds"`
	Status            string   `json:"status"`
	Author            string   `json:"author"`
}

func New(service *eduservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.POST("", h.create)
	rg.GET("/:id", h.detail)
	rg.PATCH("/:id", h.update)
	rg.POST("/:id/send", h.send)
	rg.GET("/:id/targets", h.targets)
}

func (h *Handler) list(c *gin.Context) {
	if !permission.HasFromContext(c, "notices:view") {
		response.Forbidden(c)
		return
	}

	classID, classErr := parseNoticeUintParam(c.Query("classId"))
	if classErr != nil {
		response.Failed(c, 400, "class id is invalid")
		return
	}

	notices, noticeErr := h.service.NoticesWithFilter(eduservice.NoticeFilter{
		ClassID:    classID,
		Status:     strings.TrimSpace(c.Query("status")),
		NoticeType: strings.TrimSpace(c.Query("noticeType")),
		Date:       strings.TrimSpace(c.Query("date")),
		DateFrom:   strings.TrimSpace(c.Query("dateFrom")),
		DateTo:     strings.TrimSpace(c.Query("dateTo")),
	})
	if noticeErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, notices)
}

func (h *Handler) create(c *gin.Context) {
	if !permission.HasFromContext(c, "notices:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindNoticePayload(c)
	if !ok {
		return
	}

	createdItem, createErr := h.service.CreateNotice(input, currentOperator(c))
	if createErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Created(c, createdItem)
}

func (h *Handler) detail(c *gin.Context) {
	if !permission.HasFromContext(c, "notices:view") {
		response.Forbidden(c)
		return
	}

	notice, found, noticeErr := h.service.Notice(c.Param("id"))
	if noticeErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Success(c, gin.H{"id": c.Param("id")})
		return
	}

	response.Success(c, notice)
}

func (h *Handler) update(c *gin.Context) {
	if !permission.HasFromContext(c, "notices:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindNoticePayload(c)
	if !ok {
		return
	}

	updatedItem, found, updateErr := h.service.UpdateNotice(c.Param("id"), input, currentOperator(c))
	if updateErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "notice not found")
		return
	}

	response.Success(c, updatedItem)
}

func (h *Handler) send(c *gin.Context) {
	if !permission.HasFromContext(c, "notices:manage") {
		response.Forbidden(c)
		return
	}

	sentItem, found, sendErr := h.service.SendNotice(c.Param("id"), currentOperator(c))
	if sendErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "notice not found")
		return
	}

	response.Success(c, sentItem)
}

func (h *Handler) targets(c *gin.Context) {
	if !permission.HasFromContext(c, "notices:view") {
		response.Forbidden(c)
		return
	}

	targets, targetErr := h.service.NoticeTargets(c.Param("id"))
	if targetErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, targets)
}

func bindNoticePayload(c *gin.Context) (eduservice.NoticePayload, bool) {
	var payload noticePayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "通知表单格式不正确")
		return eduservice.NoticePayload{}, false
	}

	input := eduservice.NoticePayload{
		Title:             strings.TrimSpace(payload.Title),
		Content:           strings.TrimSpace(payload.Content),
		Category:          strings.TrimSpace(payload.Category),
		TargetScope:       strings.TrimSpace(payload.TargetScope),
		RelatedClassID:    payload.RelatedClassID,
		RelatedScheduleID: payload.RelatedScheduleID,
		StudentIDs:        deduplicateNoticeStudentIDs(payload.StudentIDs),
		Status:            strings.TrimSpace(payload.Status),
		Author:            strings.TrimSpace(payload.Author),
	}

	validationMessage := validateNoticePayload(input)
	if validationMessage != "" {
		response.Failed(c, 400, validationMessage)
		return eduservice.NoticePayload{}, false
	}

	return input, true
}

func validateNoticePayload(input eduservice.NoticePayload) string {
	if input.Title == "" {
		return "通知标题不能为空"
	}
	if input.Content == "" {
		return "通知内容不能为空"
	}
	if input.Category == "" {
		return "通知分类不能为空"
	}
	if input.TargetScope == "" {
		return "通知范围不能为空"
	}
	if input.Author == "" {
		return "发起人不能为空"
	}

	switch input.Status {
	case "草稿", "待发送", "已发送":
	default:
		return "通知状态不正确"
	}

	return ""
}

func parseNoticeUintParam(rawValue string) (uint64, error) {
	trimmedValue := strings.TrimSpace(rawValue)
	if trimmedValue == "" {
		return 0, nil
	}

	return strconv.ParseUint(trimmedValue, 10, 64)
}

func deduplicateNoticeStudentIDs(source []uint64) []uint64 {
	if len(source) == 0 {
		return []uint64{}
	}

	items := make([]uint64, 0, len(source))
	seenIDs := make(map[uint64]struct{}, len(source))
	for _, studentID := range source {
		if studentID == 0 {
			continue
		}
		if _, exists := seenIDs[studentID]; exists {
			continue
		}

		seenIDs[studentID] = struct{}{}
		items = append(items, studentID)
	}

	return items
}

func currentOperator(c *gin.Context) eduservice.Operator {
	return eduservice.Operator{
		UserID:      c.GetUint64("current_user_id"),
		DisplayName: c.GetString("current_user_name"),
	}
}
