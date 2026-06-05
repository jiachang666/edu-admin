package handler

import (
	"edu-admin/internal/app/permission"
	"edu-admin/internal/app/response"
	eduservice "edu-admin/internal/modules/edu/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *eduservice.Service
}

type coursePayload struct {
	Name                  string `json:"name"`
	Category              string `json:"category"`
	Description           string `json:"description"`
	AgeRange              string `json:"ageRange"`
	LessonDurationMinutes int    `json:"lessonDurationMinutes"`
	TotalLessons          int    `json:"totalLessons"`
	DeliveryType          string `json:"deliveryType"`
	Status                string `json:"status"`
}

func New(service *eduservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.POST("", h.create)
	rg.GET("/options", h.options)
	rg.GET("/:id", h.detail)
	rg.PATCH("/:id", h.update)
}

func (h *Handler) list(c *gin.Context) {
	if !permission.HasFromContext(c, "courses:view") {
		response.Forbidden(c)
		return
	}

	filter := eduservice.CourseFilter{
		Keyword:  c.Query("keyword"),
		Category: c.Query("category"),
		Status:   c.Query("status"),
	}

	courses, courseErr := h.service.Courses(filter)
	if courseErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, courses)
}

func (h *Handler) create(c *gin.Context) {
	if !permission.HasFromContext(c, "courses:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindCoursePayload(c)
	if !ok {
		return
	}

	createdItem, createErr := h.service.CreateCourse(input)
	if createErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Created(c, createdItem)
}

func (h *Handler) options(c *gin.Context) {
	if !permission.HasFromContext(c, "courses:view") {
		response.Forbidden(c)
		return
	}

	options, optionErr := h.service.CourseOptions()
	if optionErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, options)
}

func (h *Handler) detail(c *gin.Context) {
	if !permission.HasFromContext(c, "courses:view") {
		response.Forbidden(c)
		return
	}

	course, found, courseErr := h.service.Course(c.Param("id"))
	if courseErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "course not found")
		return
	}

	response.Success(c, course)
}

func (h *Handler) update(c *gin.Context) {
	if !permission.HasFromContext(c, "courses:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindCoursePayload(c)
	if !ok {
		return
	}

	updatedItem, found, updateErr := h.service.UpdateCourse(c.Param("id"), input)
	if updateErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "course not found")
		return
	}

	response.Success(c, updatedItem)
}

func bindCoursePayload(c *gin.Context) (eduservice.CoursePayload, bool) {
	var payload coursePayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "课程表单格式不正确")
		return eduservice.CoursePayload{}, false
	}

	input := eduservice.CoursePayload{
		Name:                  strings.TrimSpace(payload.Name),
		Category:              strings.TrimSpace(payload.Category),
		Description:           strings.TrimSpace(payload.Description),
		AgeRange:              strings.TrimSpace(payload.AgeRange),
		LessonDurationMinutes: payload.LessonDurationMinutes,
		TotalLessons:          payload.TotalLessons,
		DeliveryType:          strings.TrimSpace(payload.DeliveryType),
		Status:                strings.TrimSpace(payload.Status),
	}

	validationMessage := validateCoursePayload(input)
	if validationMessage != "" {
		response.Failed(c, 400, validationMessage)
		return eduservice.CoursePayload{}, false
	}

	return input, true
}

func validateCoursePayload(input eduservice.CoursePayload) string {
	if input.Name == "" {
		return "课程名称不能为空"
	}
	if input.Category == "" {
		return "课程分类不能为空"
	}
	if input.LessonDurationMinutes <= 0 {
		return "单次时长必须大于 0"
	}
	if input.TotalLessons <= 0 {
		return "建议总课时必须大于 0"
	}
	if input.DeliveryType == "" {
		return "授课方式不能为空"
	}

	switch input.Status {
	case "启用", "停用":
	default:
		return "课程状态不正确"
	}

	return ""
}
