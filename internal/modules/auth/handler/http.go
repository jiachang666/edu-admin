package handler

import (
	"edu-admin/internal/app/response"
	authservice "edu-admin/internal/modules/auth/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *authservice.Service
}

func New(service *authservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", h.login)
	rg.POST("/logout", h.logout)
	rg.GET("/me", h.me)
	rg.GET("/me/permissions", h.permissions)
}

func (h *Handler) login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		response.Failed(c, 400, "登录表单格式不正确")
		return
	}

	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if username == "" || password == "" {
		response.Failed(c, 400, "账号和密码不能为空")
		return
	}

	result, loginErr := h.service.Login(username, password)
	if loginErr != nil {
		switch loginErr {
		case authservice.ErrInvalidCredentials:
			response.Failed(c, 401, "账号或密码不正确")
			return
		case authservice.ErrUserDisabled:
			response.Failed(c, 403, "当前账号已停用")
			return
		default:
			response.InternalServerError(c)
			return
		}
	}

	response.Success(c, result)
}

func (h *Handler) logout(c *gin.Context) {
	response.Success(c, gin.H{"loggedOut": true})
}

func (h *Handler) me(c *gin.Context) {
	currentUserID := c.GetUint64("current_user_id")
	result, found, meErr := h.service.Me(currentUserID)
	if meErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Unauthorized(c)
		return
	}

	response.Success(c, result)
}

func (h *Handler) permissions(c *gin.Context) {
	currentUserID := c.GetUint64("current_user_id")
	result, found, permissionErr := h.service.Permissions(currentUserID)
	if permissionErr != nil {
		response.InternalServerError(c)
		return
	}
	if !found {
		response.Unauthorized(c)
		return
	}

	response.Success(c, result)
}
