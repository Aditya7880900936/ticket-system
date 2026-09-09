package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
}

func Load() *Config {
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		dbPort = 5432
	}

	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "ticket_system"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}