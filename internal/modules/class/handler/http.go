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

type classPayload struct {
	Name           string `json:"name"`
	CourseID       uint64 `json:"courseId"`
	TeacherID      uint64 `json:"teacherId"`
	Campus         string `json:"campus"`
	Capacity       int    `json:"capacity"`
	WeeklySchedule string `json:"weeklySchedule"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
	Status         string `json:"status"`
	Remark         string `json:"remark"`
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
	if !permission.HasFromContext(c, "classes:view") {
		response.Forbidden(c)
		return
	}

	courseID, courseErr := parseClassUintParam(c.Query("courseId"))
	if courseErr != nil {
		response.Failed(c, 400, "course id is invalid")
		return
	}

	teacherID, teacherErr := parseClassUintParam(c.Query("teacherId"))
	if teacherErr != nil {
		response.Failed(c, 400, "teacher id is invalid")
		return
	}

	classes, classErr := h.service.ClassesWithFilter(eduservice.ClassFilter{
		Keyword:   strings.TrimSpace(c.Query("keyword")),
		Status:    strings.TrimSpace(c.Query("status")),
		CourseID:  courseID,
		TeacherID: teacherID,
	})
	if classErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, classes)
}

func (h *Handler) create(c *gin.Context) {
	if !permission.HasFromContext(c, "classes:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindClassPayload(c)
	if !ok {
		return
	}

	createdItem, createErr := h.service.CreateClass(input)
	if createErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Created(c, createdItem)
}

func (h *Handler) detail(c *gin.Context) {
	if !permission.HasFromContext(c, "classes:view") {
		response.Forbidden(c)
		return
	}

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
	if !permission.HasFromContext(c, "classes:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindClassPayload(c)
	if !ok {
		return
	}

	updatedItem, found, updateErr := h.service.UpdateClass(c.Param("id"), input)
	if updateErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "class not found")
		return
	}

	response.Success(c, updatedItem)
}

func (h *Handler) students(c *gin.Context) {
	if !permission.HasFromContext(c, "classes:view") {
		response.Forbidden(c)
		return
	}

	students, studentErr := h.service.ClassStudents(c.Param("id"))
	if studentErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, students)
}

func (h *Handler) addStudents(c *gin.Context) {
	if !permission.HasFromContext(c, "classes:manage") {
		response.Forbidden(c)
		return
	}

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
	if !permission.HasFromContext(c, "classes:manage") {
		response.Forbidden(c)
		return
	}

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
	if !permission.HasFromContext(c, "classes:view") {
		response.Forbidden(c)
		return
	}

	schedules, scheduleErr := h.service.UpcomingSchedules(c.Param("id"))
	if scheduleErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, schedules)
}

func bindClassPayload(c *gin.Context) (eduservice.ClassPayload, bool) {
	var payload classPayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "班级表单格式不正确")
		return eduservice.ClassPayload{}, false
	}

	input := eduservice.ClassPayload{
		Name:           strings.TrimSpace(payload.Name),
		CourseID:       payload.CourseID,
		TeacherID:      payload.TeacherID,
		Campus:         strings.TrimSpace(payload.Campus),
		Capacity:       payload.Capacity,
		WeeklySchedule: strings.TrimSpace(payload.WeeklySchedule),
		StartDate:      strings.TrimSpace(payload.StartDate),
		EndDate:        strings.TrimSpace(payload.EndDate),
		Status:         strings.TrimSpace(payload.Status),
		Remark:         strings.TrimSpace(payload.Remark),
	}

	validationMessage := validateClassPayload(input)
	if validationMessage != "" {
		response.Failed(c, 400, validationMessage)
		return eduservice.ClassPayload{}, false
	}

	return input, true
}

func parseClassUintParam(rawValue string) (uint64, error) {
	trimmedValue := strings.TrimSpace(rawValue)
	if trimmedValue == "" {
		return 0, nil
	}

	return strconv.ParseUint(trimmedValue, 10, 64)
}

func validateClassPayload(input eduservice.ClassPayload) string {
	if input.Name == "" {
		return "班级名称不能为空"
	}
	if input.CourseID == 0 {
		return "请选择课程"
	}
	if input.TeacherID == 0 {
		return "请选择主讲老师"
	}
	if input.Campus == "" {
		return "所属校区不能为空"
	}
	if input.Capacity <= 0 {
		return "班级容量必须大于 0"
	}
	if input.WeeklySchedule == "" {
		return "固定排课不能为空"
	}

	switch input.Status {
	case "开班中", "待满班", "已结课", "已停班":
	default:
		return "班级状态不正确"
	}

	return ""
}
