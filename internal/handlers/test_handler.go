package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/historian/backend/internal/middleware"
	"github.com/historian/backend/internal/models"
	"github.com/historian/backend/internal/services"
)

// TestHandler handles tests.
type TestHandler struct {
	svc *services.TestService
}

// NewTestHandler constructor.
func NewTestHandler(s *services.TestService) *TestHandler { return &TestHandler{svc: s} }

// List handler.
func (h *TestHandler) List(w http.ResponseWriter, r *http.Request) {
	tests, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// hide correct answers
	for i := range tests {
		for j := range tests[i].Questions {
			tests[i].Questions[j].Correct = 0
		}
	}
	writeJSON(w, http.StatusOK, tests)
}

// Get handler.
func (h *TestHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// Submit handler.
func (h *TestHandler) Submit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	u, _ := r.Context().Value(middleware.UserKey).(*models.User)
	var req services.SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	res, err := h.svc.Submit(r.Context(), id, u.ID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// MyResults handler.
func (h *TestHandler) MyResults(w http.ResponseWriter, r *http.Request) {
	u, _ := r.Context().Value(middleware.UserKey).(*models.User)
	res, err := h.svc.MyResults(r.Context(), u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, res)
}
