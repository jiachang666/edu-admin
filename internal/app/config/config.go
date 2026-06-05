package config

import (
	"os"
	"strings"
)

type Config struct {
	AppName       string
	AppEnv        string
	HTTPAddr      string
	MySQLDSN      string
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string
	MySQLAutoSeed bool
	DevAuthToken  string
}

func Load() *Config {
	loadDotEnv(".env")

	return &Config{
		AppName:       getenv("APP_NAME", "Edu Admin"),
		AppEnv:        getenv("APP_ENV", "local"),
		HTTPAddr:      getenv("HTTP_ADDR", ":8080"),
		MySQLDSN:      getenv("MYSQL_DSN", ""),
		MySQLHost:     getenv("MYSQL_HOST", "127.0.0.1"),
		MySQLPort:     getenv("MYSQL_PORT", "3306"),
		MySQLUser:     getenv("MYSQL_USER", "root"),
		MySQLPassword: getenv("MYSQL_PASSWORD", ""),
		MySQLDatabase: getenv("MYSQL_DATABASE", "edu_admin"),
		MySQLAutoSeed: getenvBool("MYSQL_AUTO_SEED", true),
		DevAuthToken:  getenv("DEV_AUTH_TOKEN", "dev-token"),
	}
}

func (c *Config) MySQLConfigured() bool {
	if strings.TrimSpace(c.MySQLDSN) != "" {
		return true
	}

	return strings.TrimSpace(c.MySQLHost) != "" &&
		strings.TrimSpace(c.MySQLUser) != "" &&
		strings.TrimSpace(c.MySQLDatabase) != ""
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getenvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return fallback
	}

	return value != "0" && value != "false" && value != "no"
}
