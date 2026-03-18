package service

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/user/access-storage-server/internal/model"
	"github.com/user/access-storage-server/internal/repository"
)

type AdminService struct {
	users  *repository.UserRepository
	shares *repository.ShareRepository
}

func NewAdminService(users *repository.UserRepository, shares *repository.ShareRepository) *AdminService {
	return &AdminService{users: users, shares: shares}
}

func (s *AdminService) IsAdmin(userID int64) (bool, error) {
	user, err := s.users.GetByID(userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, fmt.Errorf("user not found")
	}
	return user.IsAdmin, nil
}

func (s *AdminService) ListUsers() ([]model.User, error) {
	return s.users.ListAll()
}

func (s *AdminService) CreateUser(req model.AdminCreateUserRequest) (*model.User, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	existing, err := s.users.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	return s.users.CreateWithAdmin(email, string(hash), req.IsAdmin)
}

func (s *AdminService) UpdateUser(id int64, req model.AdminUpdateUserRequest) error {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" {
		return fmt.Errorf("email is required")
	}

	existing, err := s.users.GetByEmail(email)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != id {
		return fmt.Errorf("email already taken")
	}

	if !req.IsAdmin {
		user, err := s.users.GetByID(id)
		if err != nil {
			return err
		}
		if user != nil && user.IsAdmin {
			count, err := s.users.CountAdmins()
			if err != nil {
				return err
			}
			if count <= 1 {
				return fmt.Errorf("cannot remove admin role: this is the last administrator")
			}
		}
	}

	var passwordHash string
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		passwordHash = string(hash)
	}

	return s.users.UpdateProfile(id, email, passwordHash, req.IsAdmin)
}

func (s *AdminService) ListUserShares(userID int64) ([]model.ProjectShare, error) {
	return s.shares.ListSharesByUser(userID)
}

func (s *AdminService) DeleteUser(id int64) error {
	user, err := s.users.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	if user.IsAdmin {
		count, err := s.users.CountAdmins()
		if err != nil {
			return err
		}
		if count <= 1 {
			return fmt.Errorf("cannot delete the last administrator")
		}
	}
	return s.users.Delete(id)
}

func (s *AdminService) ShareProject(projectID, userID, sharedBy int64) error {
	return s.shares.Share(projectID, userID, sharedBy)
}

func (s *AdminService) UnshareProject(projectID, userID int64) error {
	return s.shares.Unshare(projectID, userID)
}
