package config

// Package config is the only backend layer that reads environment variables.
// The rest of the application receives a validated Config value, which makes
// dependencies explicit and keeps os.Getenv calls out of business logic.

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains every runtime setting used by the API. Values come from the
// process environment (for example apps/api/.env locally or Docker Compose on
// a server), not from data sent by the frontend.
type Config struct {
	Address         string
	Environment     string
	DatabaseURL     string
	MigrateOnStart  bool
	CookieName      string
	CookieSecure    bool
	SessionTTL      time.Duration
	OTPExpiry       time.Duration
	OTPCooldown     time.Duration
	OTPMaxAttempts  int
	OTPPepper       string
	EmailDelivery   string
	ExposeDevCodes  bool
	SMTPHost        string
	SMTPPort        string
	SMTPUser        string
	SMTPPassword    string
	SMTPFrom        string
	AdminEmails     map[string]struct{}
	ShutdownTimeout time.Duration
}

// Load reads environment variables, supplies safe development defaults, and
// rejects combinations that would make the API unsafe or unable to start.
// It is called exactly once from cmd/api/main.go.
func Load() (Config, error) {
	cfg := Config{
		Address:         env("API_ADDR", ":8080"),
		Environment:     env("APP_ENV", "development"),
		DatabaseURL:     strings.TrimSpace(os.Getenv("DATABASE_URL")),
		MigrateOnStart:  envBool("MIGRATE_ON_START", true),
		CookieName:      env("SESSION_COOKIE_NAME", "eshche_est_session"),
		CookieSecure:    envBool("SESSION_COOKIE_SECURE", false),
		SessionTTL:      envDuration("SESSION_TTL", 30*24*time.Hour),
		OTPExpiry:       envDuration("OTP_EXPIRY", 10*time.Minute),
		OTPCooldown:     envDuration("OTP_COOLDOWN", time.Minute),
		OTPMaxAttempts:  envInt("OTP_MAX_ATTEMPTS", 5),
		OTPPepper:       os.Getenv("OTP_PEPPER"),
		EmailDelivery:   strings.ToLower(env("EMAIL_DELIVERY", "development")),
		ExposeDevCodes:  envBool("EXPOSE_DEV_CODES", false),
		SMTPHost:        strings.TrimSpace(os.Getenv("SMTP_HOST")),
		SMTPPort:        env("SMTP_PORT", "587"),
		SMTPUser:        strings.TrimSpace(os.Getenv("SMTP_USER")),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:        strings.TrimSpace(os.Getenv("SMTP_FROM")),
		AdminEmails:     parseEmailSet(os.Getenv("ADMIN_EMAILS")),
		ShutdownTimeout: envDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}

	// Validate after collecting all values. An error returned here reaches
	// main.go, which logs it and exits before opening the HTTP port.
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if len(cfg.OTPPepper) < 32 {
		return Config{}, errors.New("OTP_PEPPER must contain at least 32 characters")
	}
	if cfg.OTPMaxAttempts < 1 || cfg.OTPMaxAttempts > 20 {
		return Config{}, errors.New("OTP_MAX_ATTEMPTS must be between 1 and 20")
	}
	if cfg.EmailDelivery != "development" && cfg.EmailDelivery != "smtp" {
		return Config{}, errors.New("EMAIL_DELIVERY must be development or smtp")
	}
	if cfg.EmailDelivery == "smtp" {
		if cfg.SMTPHost == "" || cfg.SMTPFrom == "" {
			return Config{}, errors.New("SMTP_HOST and SMTP_FROM are required in smtp mode")
		}
		if cfg.ExposeDevCodes {
			return Config{}, errors.New("EXPOSE_DEV_CODES cannot be enabled in smtp mode")
		}
	}

	return cfg, nil
}

// IsAdmin reports whether a normalized email is in the comma-separated
// ADMIN_EMAILS setting. It is used only while registering a new account.
func (c Config) IsAdmin(email string) bool {
	_, ok := c.AdminEmails[strings.ToLower(strings.TrimSpace(email))]
	return ok
}

// env returns a trimmed string or a fallback when the variable is absent.
func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

// envBool deliberately uses the fallback for an absent or malformed value.
func envBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

// envInt is the integer equivalent of envBool.
func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

// envDuration accepts Go duration strings such as "30s", "10m", or "24h".
func envDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

// parseEmailSet turns "a@example.com,b@example.com" into a set. The empty
// struct occupies no data and map lookup gives efficient membership checks.
func parseEmailSet(raw string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, item := range strings.Split(raw, ",") {
		if email := strings.ToLower(strings.TrimSpace(item)); email != "" {
			result[email] = struct{}{}
		}
	}
	return result
}

// String returns only non-secret fields so Config can safely appear in logs.
// Do not add DATABASE_URL, OTP_PEPPER, or SMTP_PASSWORD here.
func (c Config) String() string {
	return fmt.Sprintf("environment=%s address=%s email_delivery=%s", c.Environment, c.Address, c.EmailDelivery)
}
