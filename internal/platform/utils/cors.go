package utils

import (
	"net/http"
	"strings"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) setup.
func CORSMiddleware(allowedOriginsStr string) func(http.Handler) http.Handler {
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		parts := strings.Split(allowedOriginsStr, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				isAllowed := false
				for _, o := range allowedOrigins {
					if o == origin || o == "*" {
						isAllowed = true
						break
					}
				}

				if isAllowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
