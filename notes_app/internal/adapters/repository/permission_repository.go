package repository

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type PermissionRepository struct {
	queries *schemas.Queries
}

func NewPermissionRepository(queries *schemas.Queries) output.PermissionRepository {
	return &PermissionRepository{queries: queries}
}

func (r *PermissionRepository) GetAll(ctx context.Context) (*[]schemas.GetAllPermissionsRow, error) {
	permissions, err := r.queries.GetAllPermissions(ctx)

	if err != nil {
		return nil, err
	}

	return &permissions, nil
}
