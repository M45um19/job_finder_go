package router

import (
	"jobfinder/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(authHandler *handlers.AuthHandler) http.Handler {

	r := chi.NewRouter()

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	return r
}
