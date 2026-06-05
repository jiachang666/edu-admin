package handler

import (
	"edu-admin/internal/app/response"
	eduservice "edu-admin/internal/modules/edu/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *eduservice.Service
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
	response.Created(c, gin.H{"id": 4})
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
	response.Success(c, gin.H{"updated": true})
}

func (h *Handler) send(c *gin.Context) {
	response.Success(c, gin.H{"sent": true})
}

func (h *Handler) targets(c *gin.Context) {
	targets, targetErr := h.service.NoticeTargets(c.Param("id"))
	if targetErr != nil {
		response.InternalServerError(c)
		return
	}

	response.Success(c, targets)
}
