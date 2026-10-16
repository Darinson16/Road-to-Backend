package repository

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type ActionRepository struct {
	queries *schemas.Queries
}

func NewActionRepository(queries *schemas.Queries) output.ActionRepository {
	return &ActionRepository{queries: queries}
}

func (r *ActionRepository) GetAll(ctx context.Context) (*[]schemas.Action, error) {
	actions, err := r.queries.GetAllActions(ctx)

	if err != nil {
		return nil, err
	}

	return &actions, nil
}

func (r *ActionRepository) GetById(ctx context.Context, id int32) (*schemas.Action, error) {
	action, err := r.queries.GetActionById(ctx, id)

	if err != nil {
		return nil, err
	}

	return &action, nil
}
