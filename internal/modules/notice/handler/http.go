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

type noticePayload struct {
	Title          string `json:"title"`
	Content        string `json:"content"`
	Category       string `json:"category"`
	TargetScope    string `json:"targetScope"`
	RelatedClassID uint64 `json:"relatedClassId"`
	Status         string `json:"status"`
	Author         string `json:"author"`
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
	notices, noticeErr := h.service.Notices()
	if noticeErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, notices)
}

func (h *Handler) create(c *gin.Context) {
	input, ok := bindNoticePayload(c)
	if !ok {
		return
	}

	createdItem, createErr := h.service.CreateNotice(input)
	if createErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Created(c, createdItem)
}

func (h *Handler) detail(c *gin.Context) {
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
	input, ok := bindNoticePayload(c)
	if !ok {
		return
	}

	updatedItem, found, updateErr := h.service.UpdateNotice(c.Param("id"), input)
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
	sentItem, found, sendErr := h.service.SendNotice(c.Param("id"))
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
		Title:          strings.TrimSpace(payload.Title),
		Content:        strings.TrimSpace(payload.Content),
		Category:       strings.TrimSpace(payload.Category),
		TargetScope:    strings.TrimSpace(payload.TargetScope),
		RelatedClassID: payload.RelatedClassID,
		Status:         strings.TrimSpace(payload.Status),
		Author:         strings.TrimSpace(payload.Author),
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
