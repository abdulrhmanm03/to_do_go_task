package config

import "os"

type Config struct {
	PORT        string
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	JWT_SECRET  string
}

func LoadConfig() *Config {
	return &Config{
		PORT:        getEnv("PORT", "8080"),
		DB_HOST:     getEnv("DB_HOST", "localhost"),
		DB_PORT:     getEnv("DB_PORT", "5432"),
		DB_USER:     getEnv("DB_USER", "postgres"),
		DB_PASSWORD: getEnv("DB_PASSWORD", "postgres"),
		DB_NAME:     getEnv("DB_NAME", "db"),
		JWT_SECRET:  getEnv("JWT_SECRET", "secret_key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
