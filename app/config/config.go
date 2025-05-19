package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	AppEnv             string
	JwtSecret          string
	RefreshSecret      string
	DbURL              string
	ResendKey          string
	GithubClientID     string
	GithubClientSecret string
}

func NewConfig() *Config {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: No .env file found")
		}
	}

	return &Config{
		AppPort:            getEnv("APP_PORT", "3000"),
		AppEnv:             getEnv("APP_ENV", "development"),
		JwtSecret:          getEnv("JWT_SECRET", ""),
		RefreshSecret:      getEnv("REFRESH_SECRET", ""),
		DbURL:              getEnv("DB_URL", ""),
		ResendKey:          getEnv("RESEND_KEY", ""),
		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
