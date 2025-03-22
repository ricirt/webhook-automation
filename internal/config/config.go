package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	RedisHost      string
	RedisPort      string
	WebhookURL     string
	WebhookAuthKey string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		DBHost:     getEnvOrDefault("DB_HOST", "postgres"),
		DBPort:     getEnvOrDefault("DB_PORT", "5432"),
		DBUser:     getEnvOrDefault("DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:     getEnvOrDefault("DB_NAME", "insider_messages"),
		RedisHost:  getEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:  getEnvOrDefault("REDIS_PORT", "6379"),
		WebhookURL: getEnvOrDefault("WEBHOOK_URL", "https://webhook.site/6039af30-b105-47eb-8e7d-883e361e5504"),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.WebhookURL == "" {
		return fmt.Errorf("WEBHOOK_URL is required")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
