package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port          string
	Env           string
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPassword    string
	DBSSLMode     string
	JWTSecret     string
	JWTExpireHours int
	UploadDir     string
	MaxFileSizeMB int64
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		Env:           getEnv("APP_ENV", "development"),
		DBHost:        getEnv("DB_HOST", "postgres"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBName:        getEnv("DB_NAME", "dentaldb"),
		DBUser:        getEnv("DB_USER", "dental"),
		DBPassword:    mustEnv("DB_PASSWORD"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		JWTSecret:     mustEnv("JWT_SECRET"),
		UploadDir:     getEnv("UPLOAD_DIR", "/uploads"),
		MaxFileSizeMB: 20,
	}
	expHours, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "8"))
	if err != nil {
		return nil, fmt.Errorf("JWT_EXPIRE_HOURS harus angka: %w", err)
	}
	cfg.JWTExpireHours = expHours
	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s timezone=Asia/Jakarta",
		c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("env variable %s wajib diisi", key))
	}
	return v
}
