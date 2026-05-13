package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/historian/backend/internal/middleware"
	"github.com/historian/backend/internal/models"
	"github.com/historian/backend/internal/services"
)

// AuthHandler HTTP handlers.
type AuthHandler struct {
	svc *services.AuthService
}

// NewAuthHandler constructor.
func NewAuthHandler(s *services.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// Register handler.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	role := models.Role(req.Role)
	if role != models.RoleTeacher && role != models.RoleStudent {
		role = models.RoleStudent
	}
	u, token, err := h.svc.Register(r.Context(), req.Email, req.Password, req.FullName, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	setSession(w, token)
	writeJSON(w, http.StatusOK, map[string]any{"user": u, "token": token})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handler.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	u, token, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	setSession(w, token)
	writeJSON(w, http.StatusOK, map[string]any{"user": u, "token": token})
}

// Logout clears cookie.
func (h *AuthHandler) Logout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.WriteHeader(http.StatusOK)
}

// Me returns current user.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	u, _ := r.Context().Value(middleware.UserKey).(*models.User)
	writeJSON(w, http.StatusOK, u)
}

func setSession(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
