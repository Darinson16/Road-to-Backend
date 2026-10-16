package service

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type PermissionService struct {
	repo output.PermissionRepository
}

func NewPermissionService(repo output.PermissionRepository) input.PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) GetAll(ctx context.Context) (*[]schemas.GetAllPermissionsRow, error) {
	return s.repo.GetAll(ctx)
}
