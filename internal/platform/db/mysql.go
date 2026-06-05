package db

import (
	"edu-admin/internal/app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(cfg *config.Config) (*gorm.DB, error) {
	if cfg.MySQLDSN == "" {
		return nil, nil
	}

	return gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
}
