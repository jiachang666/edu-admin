package bootstrap

import (
	"log"

	"edu-admin/internal/app/config"
	"edu-admin/internal/app/router"
	eduservice "edu-admin/internal/modules/edu/service"
	"edu-admin/internal/platform/authz"
	"edu-admin/internal/platform/db"
	xlogger "edu-admin/internal/platform/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Application struct {
	Config   *config.Config
	Logger   *log.Logger
	Engine   *gin.Engine
	DB       *gorm.DB
	Enforcer authz.Enforcer
}

func NewApplication() (*Application, error) {
	cfg := config.Load()
	logger := xlogger.New()

	dbConn, err := db.Open(cfg)
	if err != nil {
		return nil, err
	}

	eduSvc := eduservice.New(dbConn)
	bootstrapErr := eduSvc.Bootstrap(cfg.MySQLAutoSeed)
	if bootstrapErr != nil {
		return nil, bootstrapErr
	}

	enforcer := authz.NewNoopEnforcer()
	engine := router.New(cfg, logger, eduSvc)

	return &Application{
		Config:   cfg,
		Logger:   logger,
		Engine:   engine,
		DB:       dbConn,
		Enforcer: enforcer,
	}, nil
}

func (a *Application) Run() error {
	a.Logger.Printf("starting %s on %s", a.Config.AppName, a.Config.HTTPAddr)
	return a.Engine.Run(a.Config.HTTPAddr)
}
