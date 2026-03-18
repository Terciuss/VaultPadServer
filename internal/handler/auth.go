package handler

import (
	"encoding/json"
	"net/http"

	"github.com/user/access-storage-server/internal/middleware"
	"github.com/user/access-storage-server/internal/model"
	"github.com/user/access-storage-server/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	user, err := h.authService.GetUser(userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, user, http.StatusOK)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req model.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.CurrentPassword == "" {
		writeError(w, "current password is required", http.StatusBadRequest)
		return
	}

	if err := h.authService.UpdateProfile(userID, req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invalid current password" {
			status = http.StatusForbidden
		} else if err.Error() == "email already taken" {
			status = http.StatusConflict
		}
		writeError(w, err.Error(), status)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	writeJSON(w, resp, http.StatusOK)
}
