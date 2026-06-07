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

type teacherPayload struct {
	Name           string `json:"name"`
	Mobile         string `json:"mobile"`
	Title          string `json:"title"`
	MainSubject    string `json:"mainSubject"`
	EmploymentType string `json:"employmentType"`
	WeeklyHours    int    `json:"weeklyHours"`
	Campus         string `json:"campus"`
	Status         string `json:"status"`
	Remark         string `json:"remark"`
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
	if !permission.HasFromContext(c, "teachers:view") {
		response.Forbidden(c)
		return
	}

	teachers, teacherErr := h.service.Teachers()
	if teacherErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, teachers)
}

func (h *Handler) create(c *gin.Context) {
	if !permission.HasFromContext(c, "teachers:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindTeacherPayload(c)
	if !ok {
		return
	}

	createdItem, createErr := h.service.CreateTeacher(input, currentOperator(c))
	if createErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Created(c, createdItem)
}

func (h *Handler) options(c *gin.Context) {
	if !permission.HasFromContext(c, "teachers:view") {
		response.Forbidden(c)
		return
	}

	options, optionErr := h.service.TeacherOptions()
	if optionErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, options)
}

func (h *Handler) detail(c *gin.Context) {
	if !permission.HasFromContext(c, "teachers:view") {
		response.Forbidden(c)
		return
	}

	teacherDetail, found, teacherErr := h.service.TeacherDetail(c.Param("id"))
	if teacherErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "teacher not found")
		return
	}

	response.Success(c, teacherDetail)
}

func (h *Handler) update(c *gin.Context) {
	if !permission.HasFromContext(c, "teachers:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindTeacherPayload(c)
	if !ok {
		return
	}

	updatedItem, found, updateErr := h.service.UpdateTeacher(c.Param("id"), input, currentOperator(c))
	if updateErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "teacher not found")
		return
	}

	response.Success(c, updatedItem)
}

func bindTeacherPayload(c *gin.Context) (eduservice.TeacherPayload, bool) {
	var payload teacherPayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "老师表单格式不正确")
		return eduservice.TeacherPayload{}, false
	}

	input := eduservice.TeacherPayload{
		Name:           strings.TrimSpace(payload.Name),
		Mobile:         strings.TrimSpace(payload.Mobile),
		Title:          strings.TrimSpace(payload.Title),
		MainSubject:    strings.TrimSpace(payload.MainSubject),
		EmploymentType: strings.TrimSpace(payload.EmploymentType),
		WeeklyHours:    payload.WeeklyHours,
		Campus:         strings.TrimSpace(payload.Campus),
		Status:         strings.TrimSpace(payload.Status),
		Remark:         strings.TrimSpace(payload.Remark),
	}

	validationMessage := validateTeacherPayload(input)
	if validationMessage != "" {
		response.Failed(c, 400, validationMessage)
		return eduservice.TeacherPayload{}, false
	}

	return input, true
}

func validateTeacherPayload(input eduservice.TeacherPayload) string {
	if input.Name == "" {
		return "老师姓名不能为空"
	}
	if input.MainSubject == "" {
		return "主教科目不能为空"
	}
	if input.EmploymentType == "" {
		return "用工类型不能为空"
	}
	if input.WeeklyHours < 0 {
		return "周课时不能小于 0"
	}
	if input.Campus == "" {
		return "所属校区不能为空"
	}

	switch input.Status {
	case "在职", "排课中", "停用":
	default:
		return "老师状态不正确"
	}

	return ""
}

func currentOperator(c *gin.Context) eduservice.Operator {
	return eduservice.Operator{
		UserID:      c.GetUint64("current_user_id"),
		DisplayName: c.GetString("current_user_name"),
	}
}
