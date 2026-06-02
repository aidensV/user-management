package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port    string
	GinMode string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret          string
	JWTExpirationHours int
}

var AppConfig *Config

func LoadConfig() error {
	// Load .env file (ignore error if not exists)
	_ = godotenv.Load()

	jwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))

	AppConfig = &Config{
		// Server
		Port:    getEnv("PORT", "8081"),
		GinMode: getEnv("GIN_MODE", "debug"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "user_management"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// JWT
		JWTSecret:          getEnv("JWT_SECRET", "default-secret-key-change-me"),
		JWTExpirationHours: jwtExpHours,
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetJWTExpirationDuration returns the JWT expiration as time.Duration
func (c *Config) GetJWTExpirationDuration() time.Duration {
	return time.Duration(c.JWTExpirationHours) * time.Hour
}
