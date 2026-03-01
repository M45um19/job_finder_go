package app

import (
	"net/http"

	"jobfinder/internal/config"
	"jobfinder/internal/database"
	"jobfinder/internal/handlers"
	"jobfinder/internal/repository"
	"jobfinder/internal/router"
	"jobfinder/internal/services"
)

type Application struct {
	Config *config.Config
}

func New() (*Application, http.Handler) {

	cfg := config.Load()

	db := database.NewPostgresPool(cfg.DBURL)

	authRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(authRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	r := router.NewRouter(authHandler)

	return &Application{Config: cfg}, r
}
