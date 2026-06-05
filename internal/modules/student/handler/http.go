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

type studentPayload struct {
	Name             string `json:"name"`
	Gender           string `json:"gender"`
	SchoolName       string `json:"schoolName"`
	GradeName        string `json:"gradeName"`
	Campus           string `json:"campus"`
	RemainingHours   int    `json:"remainingHours"`
	Status           string `json:"status"`
	Remark           string `json:"remark"`
	GuardianName     string `json:"guardianName"`
	GuardianMobile   string `json:"guardianMobile"`
	GuardianRelation string `json:"guardianRelation"`
}

type guardianPayload struct {
	Name      string `json:"name"`
	Relation  string `json:"relation"`
	Mobile    string `json:"mobile"`
	IsPrimary bool   `json:"isPrimary"`
}

func New(service *eduservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.POST("", h.create)
	rg.GET("/:id", h.detail)
	rg.PATCH("/:id", h.update)
	rg.GET("/:id/classes", h.classes)
	rg.POST("/:id/guardians", h.createGuardian)
	rg.PATCH("/:id/guardians/:guardianId", h.updateGuardian)
	rg.DELETE("/:id/guardians/:guardianId", h.deleteGuardian)
}

func (h *Handler) list(c *gin.Context) {
	if !permission.HasFromContext(c, "students:view") {
		response.Forbidden(c)
		return
	}

	students, studentErr := h.service.Students()
	if studentErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, students)
}

func (h *Handler) create(c *gin.Context) {
	if !permission.HasFromContext(c, "students:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindStudentPayload(c)
	if !ok {
		return
	}

	createdItem, createErr := h.service.CreateStudent(input)
	if createErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Created(c, createdItem)
}

func (h *Handler) detail(c *gin.Context) {
	if !permission.HasFromContext(c, "students:view") {
		response.Forbidden(c)
		return
	}

	studentDetail, found, studentErr := h.service.StudentDetail(c.Param("id"))
	if studentErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "student not found")
		return
	}

	response.Success(c, studentDetail)
}

func (h *Handler) update(c *gin.Context) {
	if !permission.HasFromContext(c, "students:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindStudentPayload(c)
	if !ok {
		return
	}

	updatedItem, found, updateErr := h.service.UpdateStudent(c.Param("id"), input)
	if updateErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "student not found")
		return
	}

	response.Success(c, updatedItem)
}

func (h *Handler) classes(c *gin.Context) {
	if !permission.HasFromContext(c, "students:view") {
		response.Forbidden(c)
		return
	}

	classes, classErr := h.service.StudentClasses(c.Param("id"))
	if classErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, classes)
}

func (h *Handler) createGuardian(c *gin.Context) {
	if !permission.HasFromContext(c, "students:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindGuardianPayload(c)
	if !ok {
		return
	}

	createdItem, found, createErr := h.service.CreateStudentGuardian(c.Param("id"), input)
	if createErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "student not found")
		return
	}

	response.Created(c, createdItem)
}

func (h *Handler) updateGuardian(c *gin.Context) {
	if !permission.HasFromContext(c, "students:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindGuardianPayload(c)
	if !ok {
		return
	}

	updatedItem, found, updateErr := h.service.UpdateStudentGuardian(c.Param("id"), c.Param("guardianId"), input)
	if updateErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "guardian not found")
		return
	}

	response.Success(c, updatedItem)
}

func (h *Handler) deleteGuardian(c *gin.Context) {
	if !permission.HasFromContext(c, "students:manage") {
		response.Forbidden(c)
		return
	}

	deleted, deleteErr := h.service.DeleteStudentGuardian(c.Param("id"), c.Param("guardianId"))
	if deleteErr != nil {
		response.InternalServerError(c)
		return
	}
	if !deleted {
		response.Failed(c, 404, "guardian not found")
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

func bindStudentPayload(c *gin.Context) (eduservice.StudentPayload, bool) {
	var payload studentPayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "学员表单格式不正确")
		return eduservice.StudentPayload{}, false
	}

	input := eduservice.StudentPayload{
		Name:             strings.TrimSpace(payload.Name),
		Gender:           strings.TrimSpace(payload.Gender),
		SchoolName:       strings.TrimSpace(payload.SchoolName),
		GradeName:        strings.TrimSpace(payload.GradeName),
		Campus:           strings.TrimSpace(payload.Campus),
		RemainingHours:   payload.RemainingHours,
		Status:           strings.TrimSpace(payload.Status),
		Remark:           strings.TrimSpace(payload.Remark),
		GuardianName:     strings.TrimSpace(payload.GuardianName),
		GuardianMobile:   strings.TrimSpace(payload.GuardianMobile),
		GuardianRelation: strings.TrimSpace(payload.GuardianRelation),
	}

	validationMessage := validateStudentPayload(input)
	if validationMessage != "" {
		response.Failed(c, 400, validationMessage)
		return eduservice.StudentPayload{}, false
	}

	return input, true
}

func bindGuardianPayload(c *gin.Context) (eduservice.StudentGuardianPayload, bool) {
	var payload guardianPayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "家长表单格式不正确")
		return eduservice.StudentGuardianPayload{}, false
	}

	input := eduservice.StudentGuardianPayload{
		Name:      strings.TrimSpace(payload.Name),
		Relation:  strings.TrimSpace(payload.Relation),
		Mobile:    strings.TrimSpace(payload.Mobile),
		IsPrimary: payload.IsPrimary,
	}

	validationMessage := validateGuardianPayload(input)
	if validationMessage != "" {
		response.Failed(c, 400, validationMessage)
		return eduservice.StudentGuardianPayload{}, false
	}

	return input, true
}

func validateStudentPayload(input eduservice.StudentPayload) string {
	if input.Name == "" {
		return "学员姓名不能为空"
	}
	if input.GradeName == "" {
		return "学员年级不能为空"
	}
	if input.Campus == "" {
		return "所属校区不能为空"
	}
	if input.RemainingHours < 0 {
		return "剩余课时不能小于 0"
	}
	if input.GuardianName == "" {
		return "至少填写一个家长联系人"
	}
	if input.GuardianMobile == "" {
		return "家长手机号不能为空"
	}

	switch input.Status {
	case "在读", "待续费", "停课", "结课":
	default:
		return "学员状态不正确"
	}

	return ""
}

func validateGuardianPayload(input eduservice.StudentGuardianPayload) string {
	if input.Name == "" {
		return "家长姓名不能为空"
	}
	if input.Mobile == "" {
		return "家长手机号不能为空"
	}

	return ""
}
