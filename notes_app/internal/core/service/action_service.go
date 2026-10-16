package service

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
	app_errors "github.com/Darinson16/notesgo/internal/core/errors"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type ActionService struct {
	repo output.ActionRepository
}

func NewActionService(repo output.ActionRepository) input.ActionService {
	return &ActionService{repo: repo}
}

func (s *ActionService) GetAll(ctx context.Context) (*[]schemas.Action, error) {
	return s.repo.GetAll(ctx)
}

func (s *ActionService) GetById(ctx context.Context, id int32) (*schemas.Action, error) {
	if id == 0 {
		return nil, app_errors.ErrIdNotProvided
	}

	return s.repo.GetById(ctx, id)
}
