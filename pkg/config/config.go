package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Namespace           string
	PrometheusRetention string
	GrafanaAdminPassword string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		Namespace:           getEnv("MONITORING_NAMESPACE", "monitoring"),
		PrometheusRetention: getEnv("PROMETHEUS_RETENTION", "15d"),
		GrafanaAdminPassword: getEnv("GRAFANA_ADMIN_PASSWORD", ""),
	}

	// Generate random password if not provided
	if config.GrafanaAdminPassword == "" {
		config.GrafanaAdminPassword = generateRandomPassword()
	}

	return config
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// generateRandomPassword generates a random password
func generateRandomPassword() string {
	// Simple implementation - in production, use crypto/rand
	return "prom-operator"
}

