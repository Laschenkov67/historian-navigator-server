package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/historian/backend/internal/models"
	"github.com/historian/backend/internal/services"
)

type ctxKey string

// UserKey context key.
const UserKey ctxKey = "user"

// CORS middleware.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Logger logs requests.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// Auth verifies JWT.
func Auth(svc *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		u, err := svc.ParseToken(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserKey, u)
		next(w, r.WithContext(ctx))
	}
}

// Role verifies user's role.
func Role(svc *services.AuthService, role string, next http.HandlerFunc) http.HandlerFunc {
	return Auth(svc, func(w http.ResponseWriter, r *http.Request) {
		u, _ := r.Context().Value(UserKey).(*models.User)
		if u == nil || string(u.Role) != role {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	})
}

func extractToken(r *http.Request) string {
	if c, err := r.Cookie("session"); err == nil {
		return c.Value
	}
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}
