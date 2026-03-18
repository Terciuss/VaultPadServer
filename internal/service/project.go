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

	own, err := s.projects.ListByUser(userID)
	if err != nil {
		return nil, err
	}
	shared, err := s.shares.ListSharedProjects(userID)
	if err != nil {
		return nil, err
	}

	seen := make(map[int64]bool)
	var result []model.Project
	for _, p := range own {
		seen[p.ID] = true
		result = append(result, p)
	}
	for _, p := range shared {
		if !seen[p.ID] {
			result = append(result, p)
		}
	}
	return result, nil
}

func (s *ProjectService) ListMeta(userID int64) ([]model.ProjectMeta, error) {
	if s.isAdmin(userID) {
		return s.projects.ListAllMeta()
	}

	own, err := s.projects.ListMetaByUser(userID)
	if err != nil {
		return nil, err
	}
	shared, err := s.shares.ListSharedProjectsMeta(userID)
	if err != nil {
		return nil, err
	}

	seen := make(map[int64]bool)
	var result []model.ProjectMeta
	for _, m := range own {
		seen[m.ID] = true
		result = append(result, m)
	}
	for _, m := range shared {
		if !seen[m.ID] {
			result = append(result, m)
		}
	}
	return result, nil
}

func (s *ProjectService) Get(id, userID int64) (*model.Project, error) {
	if s.isAdmin(userID) {
		p, err := s.projects.GetByIDRaw(id)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("project not found")
		}
		return p, nil
	}

	p, err := s.projects.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if p != nil {
		return p, nil
	}

	hasAccess, err := s.shares.HasAccess(id, userID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, fmt.Errorf("project not found")
	}

	p, err = s.projects.GetByIDRaw(id)
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
		UserID:           userID,
		EncryptedName:    req.EncryptedName,
		EncryptedContent: req.EncryptedContent,
		SortOrder:        req.SortOrder,
	}
	return s.projects.Create(p)
}

func (s *ProjectService) Update(id, userID int64, req model.UpdateProjectRequest) error {
	if s.isAdmin(userID) {
		existing, err := s.projects.GetByIDRaw(id)
		if err != nil {
			return err
		}
		if existing == nil {
			return fmt.Errorf("project not found")
		}
		existing.EncryptedName = req.EncryptedName
		existing.EncryptedContent = req.EncryptedContent
		existing.SortOrder = req.SortOrder
		return s.projects.UpdateRaw(existing)
	}

	existing, err := s.projects.GetByID(id, userID)
	if err != nil {
		return err
	}
	if existing != nil {
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

	existing, err = s.projects.GetByIDRaw(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("project not found")
	}
	existing.EncryptedName = req.EncryptedName
	existing.EncryptedContent = req.EncryptedContent
	existing.SortOrder = req.SortOrder
	return s.projects.UpdateRaw(existing)
}

func (s *ProjectService) Delete(id, userID int64) error {
	if s.isAdmin(userID) {
		return s.projects.DeleteRaw(id)
	}
	return s.projects.Delete(id, userID)
}
