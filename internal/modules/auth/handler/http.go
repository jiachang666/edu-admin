package handler

import (
	"edu-admin/internal/app/response"
	authservice "edu-admin/internal/modules/auth/service"

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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	response.Success(c, h.service.Login(req.Username))
}

func (h *Handler) logout(c *gin.Context) {
	response.Success(c, gin.H{"loggedOut": true})
}

func (h *Handler) me(c *gin.Context) {
	response.Success(c, h.service.Me())
}

func (h *Handler) permissions(c *gin.Context) {
	response.Success(c, gin.H{"permissions": h.service.Permissions()})
}
