package repository

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type RoleRepository struct {
	queries *schemas.Queries
}

func NewRoleRepository(queries *schemas.Queries) output.RoleRepository {
	return &RoleRepository{queries: queries}
}

func (r *RoleRepository) GetAll(ctx context.Context) (*[]schemas.Role, error) {
	Roles, err := r.queries.GetAllRoles(ctx)

	if err != nil {
		return nil, err
	}

	return &Roles, nil
}

func (r *RoleRepository) GetById(ctx context.Context, id int32) (*schemas.Role, error) {
	Role, err := r.queries.GetRoleById(ctx, id)

	if err != nil {
		return nil, err
	}

	return &Role, nil
}
