package handler

import (
	"edu-admin/internal/app/response"
	"edu-admin/internal/modules/demo"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func New() *Handler { return &Handler{} }

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.list)
	rg.POST("", h.create)
	rg.GET("/:id", h.detail)
	rg.PATCH("/:id", h.update)
	rg.POST("/:id/send", h.send)
	rg.GET("/:id/targets", h.targets)
}

func (h *Handler) list(c *gin.Context) {
	response.Paginated(c, demo.Notices())
}

func (h *Handler) create(c *gin.Context) {
	response.Created(c, gin.H{"id": 4})
}

func (h *Handler) detail(c *gin.Context) {
	notice, found := demo.FindNotice(c.Param("id"))
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
	response.Success(c, demo.NoticeTargets(c.Param("id")))
}
