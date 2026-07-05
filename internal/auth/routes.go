package auth

import "github.com/go-chi/chi/v5"

func (a *AuthHandler) RegisterRoutes(r chi.Router, authMiddleware *AuthMiddleware) {
	r.Post("/register", a.Register)
	r.Post("/login", a.Login)
	r.Post("/refresh", a.Refresh)
	r.Get("/profile/{id}", a.GetProfile)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)
		r.Put("/profile", a.UpdateProfile)
		r.Put("/profile/photo", a.UpdateProfilePhoto)
	})
}
