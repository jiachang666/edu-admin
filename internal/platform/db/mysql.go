package db

import (
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	"edu-admin/internal/app/config"

	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(cfg *config.Config) (*gorm.DB, error) {
	if !cfg.MySQLConfigured() {
		return nil, nil
	}

	if strings.TrimSpace(cfg.MySQLDSN) == "" {
		ensureErr := ensureDatabase(cfg)
		if ensureErr != nil {
			return nil, ensureErr
		}
	}

	dsn := databaseDSN(cfg)
	dbConn, openErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if openErr != nil {
		return nil, openErr
	}

	sqlDB, dbErr := dbConn.DB()
	if dbErr != nil {
		return nil, dbErr
	}

	pingErr := sqlDB.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return dbConn, nil
}

func ensureDatabase(cfg *config.Config) error {
	serverConfig := driver.Config{
		User:                 cfg.MySQLUser,
		Passwd:               cfg.MySQLPassword,
		Net:                  "tcp",
		Addr:                 net.JoinHostPort(cfg.MySQLHost, cfg.MySQLPort),
		Params:               map[string]string{"charset": "utf8mb4"},
		Loc:                  time.Local,
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	serverDSN := serverConfig.FormatDSN()

	serverDB, openErr := sql.Open("mysql", serverDSN)
	if openErr != nil {
		return openErr
	}
	defer serverDB.Close()

	pingErr := serverDB.Ping()
	if pingErr != nil {
		return pingErr
	}

	createSQL := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci",
		strings.ReplaceAll(cfg.MySQLDatabase, "`", "``"),
	)

	_, execErr := serverDB.Exec(createSQL)
	return execErr
}

func databaseDSN(cfg *config.Config) string {
	if strings.TrimSpace(cfg.MySQLDSN) != "" {
		return cfg.MySQLDSN
	}

	databaseConfig := driver.Config{
		User:                 cfg.MySQLUser,
		Passwd:               cfg.MySQLPassword,
		Net:                  "tcp",
		Addr:                 net.JoinHostPort(cfg.MySQLHost, cfg.MySQLPort),
		DBName:               cfg.MySQLDatabase,
		Params:               map[string]string{"charset": "utf8mb4"},
		Loc:                  time.Local,
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	return databaseConfig.FormatDSN()
}
