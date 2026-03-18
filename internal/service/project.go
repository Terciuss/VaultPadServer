package service

import (
	"fmt"

	"github.com/user/access-storage-server/internal/model"
	"github.com/user/access-storage-server/internal/repository"
)

type ProjectService struct {
	projects *repository.ProjectRepository
	shares   *repository.ShareRepository
	users    *repository.UserRepository
}

func NewProjectService(projects *repository.ProjectRepository, shares *repository.ShareRepository, users *repository.UserRepository) *ProjectService {
	return &ProjectService{projects: projects, shares: shares, users: users}
}

func (s *ProjectService) isAdmin(userID int64) bool {
	u, err := s.users.GetByID(userID)
	if err != nil || u == nil {
		return false
	}
	return u.IsAdmin
}

func (s *ProjectService) List(userID int64) ([]model.Project, error) {
	if s.isAdmin(userID) {
		return s.projects.ListAll()
	}
	return s.shares.ListSharedProjects(userID)
}

func (s *ProjectService) ListMeta(userID int64) ([]model.ProjectMeta, error) {
	if s.isAdmin(userID) {
		return s.projects.ListAllMeta()
	}
	return s.shares.ListSharedProjectsMeta(userID)
}

func (s *ProjectService) Get(id, userID int64) (*model.Project, error) {
	if s.isAdmin(userID) {
		p, err := s.projects.GetByID(id)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("project not found")
		}
		return p, nil
	}

	hasAccess, err := s.shares.HasAccess(id, userID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, fmt.Errorf("project not found")
	}

	p, err := s.projects.GetByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("project not found")
	}
	return p, nil
}

func (s *ProjectService) Create(userID int64, req model.CreateProjectRequest) (*model.Project, error) {
	p := &model.Project{
		EncryptedName:    req.EncryptedName,
		EncryptedContent: req.EncryptedContent,
		SortOrder:        req.SortOrder,
	}
	created, err := s.projects.Create(p)
	if err != nil {
		return nil, err
	}

	if err := s.shares.Share(created.ID, userID, userID); err != nil {
		return nil, fmt.Errorf("auto-share for creator: %w", err)
	}

	return created, nil
}

func (s *ProjectService) Update(id, userID int64, req model.UpdateProjectRequest) error {
	if s.isAdmin(userID) {
		existing, err := s.projects.GetByID(id)
		if err != nil {
			return err
		}
		if existing == nil {
			return fmt.Errorf("project not found")
		}
		existing.EncryptedName = req.EncryptedName
		existing.EncryptedContent = req.EncryptedContent
		existing.SortOrder = req.SortOrder
		return s.projects.Update(existing)
	}

	hasAccess, err := s.shares.HasAccess(id, userID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return fmt.Errorf("project not found")
	}

	existing, err := s.projects.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("project not found")
	}
	existing.EncryptedName = req.EncryptedName
	existing.EncryptedContent = req.EncryptedContent
	existing.SortOrder = req.SortOrder
	return s.projects.Update(existing)
}

func (s *ProjectService) Delete(id, userID int64) error {
	if !s.isAdmin(userID) {
		return fmt.Errorf("forbidden: only admins can delete projects")
	}
	return s.projects.Delete(id)
}
