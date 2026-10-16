package service

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/core/domain"
	app_errors "github.com/Darinson16/notesgo/internal/core/errors"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
	"github.com/Darinson16/notesgo/utils"
)

type AuthService struct {
	repo output.AuthRepository
}

func NewAuthService(repo output.AuthRepository) input.AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) SignUp(ctx context.Context, user *schemas.User) (int64, error) {
	if user == nil || user.Name == "" || user.Email == "" || user.PasswordHash == "" {
		return 0, app_errors.ErrDataIsNil
	}

	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return 0, err
	}

	user.PasswordHash = hashedPassword
	return s.repo.SignUp(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, credentials *domain.Login) (int32, error) {
	if credentials.Email == "" || credentials.Password == "" {
		return 0, app_errors.ErrDataIsNil
	}

	user, err := s.repo.Login(ctx, credentials.Email)
	if err != nil {
		return 0, app_errors.ErrInvalidCredentials
	}

	if !utils.CheckPasswordHash(credentials.Password, user.PasswordHash) {
		return 0, app_errors.ErrInvalidCredentials
	}

	return user.ID, nil
}
