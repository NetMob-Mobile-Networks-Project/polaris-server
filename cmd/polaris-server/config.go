package main

import (
	"os"
)

type Config struct {
	Port string

	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string

	APIVersion string
	APIPrefix  string

	JWTSecret string
}

func NewConfig() *Config {
	return &Config{
		Port:       getEnv("PORT", "8000"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPass:     getEnv("DB_PASS", "secret"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "polaris"),
		APIVersion: getEnv("API_VERSION", "v1"),
		APIPrefix:  getEnv("API_PREFIX", "/api"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-here"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
