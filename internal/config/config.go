package config

import (
	"log"
	"os"

	"github.com/joho/gotdotenv"
)

type Config struct {
	port string
	DBURL string
	JWTSecret string
}

func Load() *Config {

	_ := gotdotenv.Load()

	cfg := &Config {
		port: getEnv("PORT", "8080"),
		DBURL: os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET")
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}