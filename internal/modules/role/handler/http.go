package handler

import (
	"errors"
	"strings"

	"edu-admin/internal/app/permission"
	"edu-admin/internal/app/response"
	eduservice "edu-admin/internal/modules/edu/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *eduservice.Service
}

type rolePayload struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type permissionPayload struct {
	Permissions []string `json:"permissions"`
}

func New(service *eduservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.POST("", h.create)
	rg.GET("/:id", h.detail)
	rg.PATCH("/:id", h.update)
	rg.PUT("/:id/permissions", h.permissions)
}

func (h *Handler) list(c *gin.Context) {
	if !hasPermission(c, "roles:view") {
		response.Forbidden(c)
		return
	}

	roles, roleErr := h.service.Roles()
	if roleErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, roles)
}

func (h *Handler) create(c *gin.Context) {
	if !hasPermission(c, "roles:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindRolePayload(c)
	if !ok {
		return
	}

	item, createErr := h.service.CreateRole(input, currentOperator(c))
	if createErr != nil {
		handleRoleServiceError(c, createErr)
		return
	}

	response.Created(c, item)
}

func (h *Handler) detail(c *gin.Context) {
	if !hasPermission(c, "roles:view") {
		response.Forbidden(c)
		return
	}

	roleItem, found, roleErr := h.service.Role(c.Param("id"))
	if roleErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "role not found")
		return
	}

	response.Success(c, gin.H{
		"role":             roleItem,
		"permissionGroups": h.service.PermissionGroups(),
	})
}

func (h *Handler) update(c *gin.Context) {
	if !hasPermission(c, "roles:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindRolePayload(c)
	if !ok {
		return
	}

	item, found, updateErr := h.service.UpdateRole(c.Param("id"), input, currentOperator(c))
	if updateErr != nil {
		handleRoleServiceError(c, updateErr)
		return
	}
	if !found {
		response.Failed(c, 404, "role not found")
		return
	}

	response.Success(c, item)
}

func (h *Handler) permissions(c *gin.Context) {
	if !hasPermission(c, "roles:manage") {
		response.Forbidden(c)
		return
	}

	var payload permissionPayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "角色权限表单格式不正确")
		return
	}

	input := eduservice.RolePermissionPayload{
		Permissions: payload.Permissions,
	}

	item, found, saveErr := h.service.SaveRolePermissions(c.Param("id"), input, currentOperator(c))
	if saveErr != nil {
		handleRoleServiceError(c, saveErr)
		return
	}
	if !found {
		response.Failed(c, 404, "role not found")
		return
	}

	response.Success(c, item)
}

func bindRolePayload(c *gin.Context) (eduservice.RolePayload, bool) {
	var payload rolePayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "角色表单格式不正确")
		return eduservice.RolePayload{}, false
	}

	input := eduservice.RolePayload{
		Name:        strings.TrimSpace(payload.Name),
		Code:        strings.TrimSpace(payload.Code),
		Description: strings.TrimSpace(payload.Description),
		Status:      strings.TrimSpace(payload.Status),
	}

	validationMessage := validateRolePayload(input)
	if validationMessage != "" {
		response.Failed(c, 400, validationMessage)
		return eduservice.RolePayload{}, false
	}

	return input, true
}

func validateRolePayload(input eduservice.RolePayload) string {
	if input.Name == "" {
		return "角色名称不能为空"
	}
	if input.Code == "" {
		return "角色编码不能为空"
	}

	switch input.Status {
	case "启用", "停用":
	default:
		return "角色状态不正确"
	}

	return ""
}

func currentOperator(c *gin.Context) eduservice.Operator {
	return eduservice.Operator{
		UserID:      c.GetUint64("current_user_id"),
		DisplayName: c.GetString("current_user_name"),
	}
}

func hasPermission(c *gin.Context, target string) bool {
	return permission.HasFromContext(c, target)
}

func handleRoleServiceError(c *gin.Context, serviceErr error) {
	switch {
	case errors.Is(serviceErr, eduservice.ErrRoleCodeAlreadyExists):
		response.Failed(c, 400, "角色编码已存在，请换一个编码")
	case errors.Is(serviceErr, eduservice.ErrInvalidPermissionCodes):
		response.Failed(c, 400, "权限项里包含系统暂不支持的内容")
	case errors.Is(serviceErr, eduservice.ErrProtectedRole):
		response.Failed(c, 400, "内置超级管理员角色暂不允许停用或改成不完整权限")
	default:
		response.InternalServerError(c)
	}
}
