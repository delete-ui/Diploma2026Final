package config

import (
	"GolangBackendDiploma26/internal/models"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func Load() (*models.Config, error) {
	_ = godotenv.Load()

	cfg := &models.Config{
		Env:       getEnv("ENV", "development"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "text"),
		HTTPServer: models.HTTPServerConfig{
			Address:      getEnv("HTTP_ADDRESS", ":8080"),
			ReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: models.DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "battery_shop"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: models.JWTConfig{
			Secret:    getEnv("JWT_SECRET", ""),
			AccessTTL: getEnvAsDuration("JWT_ACCESS_TTL", 15*time.Minute),
		},
		SMTP: models.SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     getEnvAsInt("SMTP_PORT", 587),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@batteryshop.ru"),
			FromName: getEnv("SMTP_FROM_NAME", "Battery Shop"),
		},
	}

	if err := Validate(cfg); err != nil {
		return nil, err
	}

	log.Printf("config loaded successfully: env=%s, addr=%s, db=%s:%s/%s",
		cfg.Env, cfg.HTTPServer.Address, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	return cfg, nil
}

func MustLoad() *models.Config {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		d, err := time.ParseDuration(val)
		if err == nil {
			return d
		}
		log.Printf("warning: invalid duration for %s=%s, using default %v", key, val, defaultVal)
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		i, err := strconv.Atoi(val)
		if err == nil {
			return i
		}
		log.Printf("warning: invalid int for %s=%s", key, val)
	}
	return defaultVal
}

func DSN(c models.DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func ValidateDB(c models.DatabaseConfig) error {
	if c.Host == "" {
		return &ValidationError{Field: "DB_HOST", Message: "cannot be empty"}
	}
	if c.Port == "" {
		return &ValidationError{Field: "DB_PORT", Message: "cannot be empty"}
	}
	portNum, err := strconv.Atoi(c.Port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return &ValidationError{Field: "DB_PORT", Message: "must be a valid port number (1-65535)"}
	}
	if c.User == "" {
		return &ValidationError{Field: "DB_USER", Message: "cannot be empty"}
	}
	if c.DBName == "" {
		return &ValidationError{Field: "DB_NAME", Message: "cannot be empty"}
	}

	validSSL := map[string]bool{"disable": true, "require": true, "verify-ca": true, "verify-full": true}
	if !validSSL[strings.ToLower(c.SSLMode)] {
		return &ValidationError{Field: "DB_SSLMODE", Message: "must be one of: disable, require, verify-ca, verify-full"}
	}

	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("config validation error: %s - %s", e.Field, e.Message)
}

func Validate(c *models.Config) error {

	if c.Env != "development" && c.Env != "production" && c.Env != "test" {
		return &ValidationError{Field: "ENV", Message: "must be development, production, or test"}
	}

	if err := ValidateDB(c.Database); err != nil {
		return err
	}

	if c.HTTPServer.ReadTimeout < 0 {
		return &ValidationError{Field: "HTTP_READ_TIMEOUT", Message: "cannot be negative"}
	}
	if c.HTTPServer.WriteTimeout < 0 {
		return &ValidationError{Field: "HTTP_WRITE_TIMEOUT", Message: "cannot be negative"}
	}

	return nil
}
