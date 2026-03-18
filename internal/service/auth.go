package service

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/user/access-storage-server/internal/middleware"
	"github.com/user/access-storage-server/internal/model"
	"github.com/user/access-storage-server/internal/repository"
)

type AuthService struct {
	users     *repository.UserRepository
	jwtSecret string
}

func NewAuthService(users *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{users: users, jwtSecret: jwtSecret}
}

func (s *AuthService) GetUser(userID int64) (*model.User, error) {
	user, err := s.users.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *AuthService) Login(req model.LoginRequest) (*model.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.users.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := middleware.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) UpdateProfile(userID int64, req model.UpdateProfileRequest) error {
	user, err := s.users.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return fmt.Errorf("invalid current password")
	}

	email := user.Email
	if req.Email != "" {
		email = strings.TrimSpace(strings.ToLower(req.Email))
		existing, err := s.users.GetByEmail(email)
		if err != nil {
			return err
		}
		if existing != nil && existing.ID != userID {
			return fmt.Errorf("email already taken")
		}
	}

	var passwordHash string
	if req.NewPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		passwordHash = string(hash)
	}

	return s.users.UpdateProfile(userID, email, passwordHash, user.IsAdmin)
}
