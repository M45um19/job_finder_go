package app

import (
	"net/http"

	"jobfinder/internal/config"
	"jobfinder/internal/database"
	"jobfinder/internal/handlers"
	"jobfinder/internal/middleware"
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

	jobRepo := repository.NewJobRepository(db)
	jobService := services.NewJobService(jobRepo)
	jobHandler := handlers.NewJobHandler(jobService)

	applicationRepo := repository.NewApplicationRepository(db)
	applicationService := services.NewApplicationService(applicationRepo)
	applicationHandler := handlers.NewApplicationHandler(applicationService)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	r := router.NewRouter(authHandler, jobHandler, applicationHandler, authMiddleware)

	return &Application{Config: cfg}, r
}
