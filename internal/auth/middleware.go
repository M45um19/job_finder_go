package auth

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"jobfinder/internal/platform/utils"
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
			utils.Error(w, http.StatusUnauthorized, "Authorization header missing")
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

		var userID int64
		if val, ok := claims["user_id"].(float64); ok {
			userID = int64(val)
		} else if val, ok := claims["user_id"].(int64); ok {
			userID = val
		} else if valStr, ok := claims["user_id"].(string); ok {
			var err error
			userID, err = strconv.ParseInt(valStr, 10, 64)
			if err != nil {
				utils.Error(w, http.StatusUnauthorized, "Invalid user_id format in token")
				return
			}
		} else {
			utils.Error(w, http.StatusUnauthorized, "Invalid user_id in token")
			return
		}
		role := claims["role"].(string)

		ctx := context.WithValue(r.Context(), UserIdKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(RoleKey).(string)
			if !ok || userRole != role {
				utils.Error(w, http.StatusForbidden, "Wrong role or unauthorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
