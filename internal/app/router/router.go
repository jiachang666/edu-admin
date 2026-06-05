package router

import (
	"log"

	"edu-admin/internal/app/config"
	"edu-admin/internal/app/middleware"
	attendancehandler "edu-admin/internal/modules/attendance/handler"
	audithandler "edu-admin/internal/modules/audit/handler"
	authhandler "edu-admin/internal/modules/auth/handler"
	authservice "edu-admin/internal/modules/auth/service"
	classhandler "edu-admin/internal/modules/class/handler"
	coursehandler "edu-admin/internal/modules/course/handler"
	dashboardhandler "edu-admin/internal/modules/dashboard/handler"
	dashboardservice "edu-admin/internal/modules/dashboard/service"
	eduservice "edu-admin/internal/modules/edu/service"
	homeworkhandler "edu-admin/internal/modules/homework/handler"
	noticehandler "edu-admin/internal/modules/notice/handler"
	rolehandler "edu-admin/internal/modules/role/handler"
	schedulehandler "edu-admin/internal/modules/schedule/handler"
	studenthandler "edu-admin/internal/modules/student/handler"
	teacherhandler "edu-admin/internal/modules/teacher/handler"
	userhandler "edu-admin/internal/modules/user/handler"

	"github.com/gin-gonic/gin"
)

func New(cfg *config.Config, logger *log.Logger, eduSvc *eduservice.Service) *gin.Engine {
	engine := gin.New()
	trustedProxyErr := engine.SetTrustedProxies(nil)
	if trustedProxyErr != nil {
		logger.Printf("set trusted proxies failed: %v", trustedProxyErr)
	}
	engine.Use(middleware.RequestID(), middleware.Logger(logger), middleware.Recovery())

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := engine.Group("/api/v1")

	authSvc := authservice.New(cfg.DevAuthToken)
	authhandler.New(authSvc).RegisterRoutes(api.Group("/auth"))

	secured := api.Group("/")
	secured.Use(middleware.RequireAuth(cfg.DevAuthToken))

	dashboardhandler.New(dashboardservice.New(eduSvc)).RegisterRoutes(secured.Group("/dashboard"))
	userhandler.New().RegisterRoutes(secured.Group("/users"))
	rolehandler.New().RegisterRoutes(secured.Group("/roles"))
	teacherhandler.New(eduSvc).RegisterRoutes(secured.Group("/teachers"))
	studenthandler.New(eduSvc).RegisterRoutes(secured.Group("/students"))
	coursehandler.New().RegisterRoutes(secured.Group("/courses"))
	classhandler.New(eduSvc).RegisterRoutes(secured.Group("/classes"))
	schedulehandler.New(eduSvc).RegisterRoutes(secured.Group("/schedules"))
	attendancehandler.New().RegisterRoutes(secured.Group("/attendance"))
	homeworkhandler.New().RegisterRoutes(secured.Group("/homeworks"))
	noticehandler.New(eduSvc).RegisterRoutes(secured.Group("/notices"))
	audithandler.New().RegisterRoutes(secured.Group("/operation-logs"))

	return engine
}
