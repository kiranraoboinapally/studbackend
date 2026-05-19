package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration loaded from environment.
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT
	JWTSecret string

	// Server
	ServerPort string
	AppEnv     string
	AppURL     string

	// SMTP
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string

	// Razorpay
	RazorpayKeyID     string
	RazorpayKeySecret string

	// Outbox
	OutboxPollInterval int // seconds
	OutboxBatchSize    int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		DBHost:     env("DB_HOST", "192.168.1.201"),
		DBPort:     env("DB_PORT", "5432"),
		DBUser:     env("DB_USER", "postgres"),
		DBPassword: env("DB_PASSWORD", "root"),
		DBName:     env("DB_NAME", "university_erp_prod10"),

		JWTSecret: env("JWT_SECRET", "mySecretKeyAs123#"),

		ServerPort: env("SERVER_PORT", "8080"),
		AppEnv:     env("APP_ENV", "development"),
		AppURL:     env("APP_URL", "http://localhost:3000"),

		SMTPHost: env("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort: env("SMTP_PORT", "587"),
		SMTPUser: env("SMTP_USER", ""),
		SMTPPass: env("SMTP_PASS", ""),

		RazorpayKeyID:     env("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret: env("RAZORPAY_KEY_SECRET", ""),

		OutboxPollInterval: envInt("OUTBOX_POLL_INTERVAL", 2),
		OutboxBatchSize:    envInt("OUTBOX_BATCH_SIZE", 50),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
