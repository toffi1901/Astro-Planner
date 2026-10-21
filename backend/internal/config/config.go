package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisAddr   string
	Pepper      string
	JWTSecret   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		println("godotenv:", err.Error())
	}
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		Pepper:      getEnv("PEPPER", "dev-pepper"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret"),
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
