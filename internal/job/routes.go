package job

import (
	"github.com/go-chi/chi/v5"
	"jobfinder/internal/application"
	"jobfinder/internal/auth"
)

func (j *JobHandler) RegisterRoutes(r chi.Router, authMiddleware *auth.AuthMiddleware, appHandler *application.ApplicationHandler) {
	r.Get("/", j.GetAllJobs)
	r.Get("/{id}", j.GetSingleJobDetails)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		r.Use(authMiddleware.RequireRole("employer"))

		r.Post("/", j.CreateJob)
		r.Get("/employer", j.GetEmployerJobs)
		r.Put("/{id}", j.UpdateJob)
		r.Delete("/{id}", j.DeleteJob)

		r.Get("/{id}/applications", appHandler.GetApplicationByJobId)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		r.Use(authMiddleware.RequireRole("employee"))

		r.Post("/{id}/apply", appHandler.CreateApplication)
	})
}
