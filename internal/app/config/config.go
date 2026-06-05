package config

import "os"

type Config struct {
	AppName      string
	AppEnv       string
	HTTPAddr     string
	MySQLDSN     string
	DevAuthToken string
}

func Load() *Config {
	return &Config{
		AppName:      getenv("APP_NAME", "Edu Admin"),
		AppEnv:       getenv("APP_ENV", "local"),
		HTTPAddr:     getenv("HTTP_ADDR", ":8080"),
		MySQLDSN:     getenv("MYSQL_DSN", ""),
		DevAuthToken: getenv("DEV_AUTH_TOKEN", "dev-token"),
	}
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
