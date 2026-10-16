package service

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
	app_errors "github.com/Darinson16/notesgo/internal/core/errors"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type RoleService struct {
	repo output.RoleRepository
}

func NewRoleService(repo output.RoleRepository) input.RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetAll(ctx context.Context) (*[]schemas.Role, error) {
	return s.repo.GetAll(ctx)
}

func (s *RoleService) GetById(ctx context.Context, id int32) (*schemas.Role, error) {
	if id == 0 {
		return nil, app_errors.ErrIdNotProvided
	}

	return s.repo.GetById(ctx, id)
}
