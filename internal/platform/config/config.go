package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DBURL               string
	JWTSecret           string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
}

func Load() *Config {
	// Attempt to load .env from root first
	if err := godotenv.Load(); err != nil {
		// Fallback to configs/.env
		_ = godotenv.Load("configs/.env")
	}

	cfg := &Config{
		Port:                getEnv("PORT", "8080"),
		DBURL:               os.Getenv("DB_URL"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
