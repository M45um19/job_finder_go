package application

import (
	"github.com/go-chi/chi/v5"
	"jobfinder/internal/auth"
)

func (a *ApplicationHandler) RegisterRoutes(r chi.Router, authMiddleware *auth.AuthMiddleware) {
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		r.Use(authMiddleware.RequireRole("employee"))

		r.Get("/", a.GetApplicationByEmployeeId)
	})
}
