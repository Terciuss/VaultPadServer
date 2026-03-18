package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/user/access-storage-server/internal/middleware"
	"github.com/user/access-storage-server/internal/model"
	"github.com/user/access-storage-server/internal/service"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	projects, err := h.projectService.List(userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if projects == nil {
		projects = []model.Project{}
	}
	writeJSON(w, projects, http.StatusOK)
}

func (h *ProjectHandler) ListMeta(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	metas, err := h.projectService.ListMeta(userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if metas == nil {
		metas = []model.ProjectMeta{}
	}
	writeJSON(w, metas, http.StatusOK)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, "invalid project id", http.StatusBadRequest)
		return
	}

	project, err := h.projectService.Get(id, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, project, http.StatusOK)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req model.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	project, err := h.projectService.Create(userID, req)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, project, http.StatusCreated)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, "invalid project id", http.StatusBadRequest)
		return
	}

	var req model.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.projectService.Update(id, userID, req); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, "invalid project id", http.StatusBadRequest)
		return
	}

	if err := h.projectService.Delete(id, userID); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}
