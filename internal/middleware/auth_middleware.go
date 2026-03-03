package middleware

import (
	"context"
	"jobfinder/internal/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIdKey contextKey = "userID"
	RoleKey   contextKey = "roke"
)

type AuthMiddleware struct {
	JwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{JwtSecret: jwtSecret}
}

func (a *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Error(w, http.StatusUnauthorized, "AuthoAuthorization header missing")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.Error(w, http.StatusUnauthorized, "Invalid Token")
			return
		}
		claims := token.Claims.(jwt.MapClaims)

		userID := claims["user_id"].(int64)
		role := claims["role"].(string)

		ctx := context.WithValue(r.Context(), UserIdKey, userID)

		ctx = context.WithValue(r.Context(), RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userRole := r.Context().Value(RoleKey)

			if userRole != role {
				utils.Error(w, http.StatusForbidden, "Wrong role")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
