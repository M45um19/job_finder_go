package router

import (
	"jobfinder/internal/handlers"
	"jobfinder/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(authHandler *handlers.AuthHandler,
	jobHandler *handlers.JobHandler,
	authMiddlewre *middleware.AuthMiddleware,
) http.Handler {

	r := chi.NewRouter()

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/api/v1/jobs", func(r chi.Router) {

		r.Get("/", jobHandler.GetAllJobs)
		r.Get("/{id}", jobHandler.GetSingleJobDetails)

		r.Group(func(r chi.Router) {
			r.Use(authMiddlewre.RequireAuth)
			r.Use(authMiddlewre.RequireRole("employer"))

			r.Post("/", jobHandler.CreateJob)
		})
	})
	return r
}
