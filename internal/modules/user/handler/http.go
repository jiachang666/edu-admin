package handler

import (
	"errors"
	"strconv"
	"strings"

	"edu-admin/internal/app/permission"
	"edu-admin/internal/app/response"
	eduservice "edu-admin/internal/modules/edu/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *eduservice.Service
}

type userPayload struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	Mobile      string `json:"mobile"`
	RoleCode    string `json:"roleCode"`
	Status      string `json:"status"`
}

func New(service *eduservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.POST("", h.create)
	rg.GET("/:id", h.detail)
	rg.PATCH("/:id", h.update)
	rg.POST("/:id/enable", h.enable)
	rg.POST("/:id/disable", h.disable)
}

func (h *Handler) list(c *gin.Context) {
	if !hasPermission(c, "users:view") {
		response.Forbidden(c)
		return
	}

	roleID, roleErr := parseUserUintParam(c.Query("roleId"))
	if roleErr != nil {
		response.Failed(c, 400, "role id is invalid")
		return
	}

	users, userErr := h.service.UsersWithFilter(eduservice.UserFilter{
		Keyword: strings.TrimSpace(c.Query("keyword")),
		Status:  strings.TrimSpace(c.Query("status")),
		RoleID:  roleID,
	})
	if userErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Paginated(c, users)
}

func (h *Handler) create(c *gin.Context) {
	if !hasPermission(c, "users:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindUserPayload(c, true)
	if !ok {
		return
	}

	operator := currentOperator(c)
	item, createErr := h.service.CreateUser(input, operator)
	if createErr != nil {
		handleUserServiceError(c, createErr)
		return
	}

	response.Created(c, item)
}

func (h *Handler) detail(c *gin.Context) {
	if !hasPermission(c, "users:view") {
		response.Forbidden(c)
		return
	}

	userItem, found, userErr := h.service.User(c.Param("id"))
	if userErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Failed(c, 404, "user not found")
		return
	}

	response.Success(c, userItem)
}

func (h *Handler) update(c *gin.Context) {
	if !hasPermission(c, "users:manage") {
		response.Forbidden(c)
		return
	}

	input, ok := bindUserPayload(c, false)
	if !ok {
		return
	}

	operator := currentOperator(c)
	item, found, updateErr := h.service.UpdateUser(c.Param("id"), input, operator)
	if updateErr != nil {
		handleUserServiceError(c, updateErr)
		return
	}
	if !found {
		response.Failed(c, 404, "user not found")
		return
	}

	response.Success(c, item)
}

func (h *Handler) enable(c *gin.Context) {
	if !hasPermission(c, "users:manage") {
		response.Forbidden(c)
		return
	}

	operator := currentOperator(c)
	item, found, enableErr := h.service.EnableUser(c.Param("id"), operator)
	if enableErr != nil {
		handleUserServiceError(c, enableErr)
		return
	}
	if !found {
		response.Failed(c, 404, "user not found")
		return
	}

	response.Success(c, item)
}

func (h *Handler) disable(c *gin.Context) {
	if !hasPermission(c, "users:manage") {
		response.Forbidden(c)
		return
	}

	operator := currentOperator(c)
	item, found, disableErr := h.service.DisableUser(c.Param("id"), operator)
	if disableErr != nil {
		handleUserServiceError(c, disableErr)
		return
	}
	if !found {
		response.Failed(c, 404, "user not found")
		return
	}

	response.Success(c, item)
}

func bindUserPayload(c *gin.Context, requirePassword bool) (eduservice.UserPayload, bool) {
	var payload userPayload
	bindErr := c.ShouldBindJSON(&payload)
	if bindErr != nil {
		response.Failed(c, 400, "账号表单格式不正确")
		return eduservice.UserPayload{}, false
	}

	input := eduservice.UserPayload{
		Username:    strings.TrimSpace(payload.Username),
		Password:    strings.TrimSpace(payload.Password),
		DisplayName: strings.TrimSpace(payload.DisplayName),
		Mobile:      strings.TrimSpace(payload.Mobile),
		RoleCode:    strings.TrimSpace(payload.RoleCode),
		Status:      strings.TrimSpace(payload.Status),
	}

	validationMessage := validateUserPayload(input, requirePassword)
	if validationMessage != "" {
		response.Failed(c, 400, validationMessage)
		return eduservice.UserPayload{}, false
	}

	return input, true
}

func validateUserPayload(input eduservice.UserPayload, requirePassword bool) string {
	if input.Username == "" {
		return "账号不能为空"
	}
	if requirePassword && input.Password == "" {
		return "密码不能为空"
	}
	if input.DisplayName == "" {
		return "姓名不能为空"
	}
	if input.RoleCode == "" {
		return "角色不能为空"
	}

	switch input.Status {
	case "启用", "停用":
	default:
		return "账号状态不正确"
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

func handleUserServiceError(c *gin.Context, serviceErr error) {
	switch {
	case errors.Is(serviceErr, eduservice.ErrUsernameAlreadyExists):
		response.Failed(c, 400, "账号已存在，请换一个登录名")
	case errors.Is(serviceErr, eduservice.ErrRoleNotFound):
		response.Failed(c, 400, "角色不存在或已停用")
	case errors.Is(serviceErr, eduservice.ErrCannotDisableCurrentUser):
		response.Failed(c, 400, "不能停用当前登录账号")
	default:
		response.InternalServerError(c)
	}
}

func parseUserUintParam(rawValue string) (uint64, error) {
	trimmedValue := strings.TrimSpace(rawValue)
	if trimmedValue == "" {
		return 0, nil
	}

	return strconv.ParseUint(trimmedValue, 10, 64)
}
