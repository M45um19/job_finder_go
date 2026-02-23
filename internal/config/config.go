package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DBURL string
	JWTSecret string
}

func Load() *Config {

	_ = godotenv.Load()

	cfg := &Config {
		Port: getEnv("PORT", "8080"),
		DBURL: os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}