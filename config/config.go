package config

import (
	"iam_services_main_v1/pkg/logger"
	"os"

	"github.com/joho/godotenv"
)

// Initialize environment variables from .env file
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		logger.LogFatal("Error loading .env file")
		return err
	}
	return nil
}

// Helper function to get an environment variable or exit if not set
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.LogFatal("Environment variable %s is not set", "key", key)
	}
	return value
}
