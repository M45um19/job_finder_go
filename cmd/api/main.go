package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"jobfinder/internal/application"
	"jobfinder/internal/auth"
	"jobfinder/internal/job"
	"jobfinder/internal/platform/cloudinary"
	"jobfinder/internal/platform/config"
	"jobfinder/internal/platform/database"
	"jobfinder/internal/platform/idgen"
)

func main() {
	// 1. Load configurations
	cfg := config.Load()

	// 2. Initialize database connection pool
	db := database.NewPostgresPool(cfg.DBURL)

	// Initialize Snowflake ID generator (Node ID 1)
	idgen.Init(1)

	// 3. Initialize Cloudinary Service (optional, logs warnings if missing)
	var cldService *cloudinary.Service
	if cfg.CloudinaryCloudName != "" && cfg.CloudinaryAPIKey != "" && cfg.CloudinaryAPISecret != "" {
		var err error
		cldService, err = cloudinary.New(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
		if err != nil {
			log.Println("WARNING: Failed to initialize Cloudinary:", err)
		} else {
			log.Println("Cloudinary initialized successfully")
		}
	} else {
		log.Println("WARNING: Cloudinary credentials missing. Profile picture upload will be disabled.")
	}

	// 4. Initialize Auth Domain
	authRepo := auth.NewUserRepository(db)
	authService := auth.NewAuthService(authRepo, cfg.JWTSecret, cldService)
	authHandler := auth.NewAuthHandler(authService)
	authMiddleware := auth.NewAuthMiddleware(cfg.JWTSecret)

	// 5. Initialize Job Domain
	jobRepo := job.NewJobRepository(db)
	jobService := job.NewJobService(jobRepo)
	jobHandler := job.NewJobHandler(jobService)

	// 6. Initialize Application Domain
	appRepo := application.NewApplicationRepository(db)
	appService := application.NewApplicationService(appRepo)
	appHandler := application.NewApplicationHandler(appService)

	// 7. Build the Router
	r := chi.NewRouter()

	r.Route("/api/v1/auth", func(r chi.Router) {
		authHandler.RegisterRoutes(r, authMiddleware)
	})

	r.Route("/api/v1/jobs", func(r chi.Router) {
		jobHandler.RegisterRoutes(r, authMiddleware, appHandler)
	})

	r.Route("/api/v1/applications", func(r chi.Router) {
		appHandler.RegisterRoutes(r, authMiddleware)
	})

	// 8. Start HTTP Server
	log.Println("server is running on port:", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}