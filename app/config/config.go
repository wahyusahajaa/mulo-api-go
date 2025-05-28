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
	DBHost             string
	DBPort             string
	DBUser             string
	DBPass             string
	DBname             string
	DBSSLMode          string
	DBSSLRootCert      string
	ResendKey          string
	GithubClientID     string
	GithubClientSecret string
}

func NewConfig() *Config {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env.development")
		if err != nil {
			log.Println("Warning: .env.development file not found")
		}
	}

	return &Config{
		AppPort:            getEnv("APP_PORT", "8081"),
		AppEnv:             getEnv("APP_ENV", "development"),
		JwtSecret:          getEnv("JWT_SECRET", ""),
		RefreshSecret:      getEnv("REFRESH_SECRET", ""),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPass:             getEnv("DB_PASS", "postgres"),
		DBname:             getEnv("DB_NAME", "postgres"),
		DBSSLMode:          getEnv("DB_SSL_MODE", "disable"),
		DBSSLRootCert:      getEnv("DB_SSL_ROOT_CERT", ""),
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
