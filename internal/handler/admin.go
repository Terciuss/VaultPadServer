package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/user/access-storage-server/internal/middleware"
	"github.com/user/access-storage-server/internal/model"
	"github.com/user/access-storage-server/internal/service"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	userID := middleware.GetUserID(r)
	isAdmin, err := h.adminService.IsAdmin(userID)
	if err != nil || !isAdmin {
		writeError(w, "admin access required", http.StatusForbidden)
		return false
	}
	return true
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	users, err := h.adminService.ListUsers()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if users == nil {
		users = []model.User{}
	}
	writeJSON(w, users, http.StatusOK)
}

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	var req model.AdminCreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.adminService.CreateUser(req)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, user, http.StatusCreated)
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req model.AdminUpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.adminService.UpdateUser(id, req); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if err := h.adminService.DeleteUser(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func (h *AdminHandler) ListUserShares(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	shares, err := h.adminService.ListUserShares(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if shares == nil {
		shares = []model.ProjectShare{}
	}
	writeJSON(w, shares, http.StatusOK)
}

func (h *AdminHandler) ShareProject(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	projectID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, "invalid project id", http.StatusBadRequest)
		return
	}

	var req model.ShareProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sharedBy := middleware.GetUserID(r)
	if err := h.adminService.ShareProject(projectID, req.UserID, sharedBy); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"}, http.StatusCreated)
}

func (h *AdminHandler) UnshareProject(w http.ResponseWriter, r *http.Request) {
	if !h.requireAdmin(w, r) {
		return
	}

	projectID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, "invalid project id", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("userId"), 10, 64)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if err := h.adminService.UnshareProject(projectID, userID); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}
