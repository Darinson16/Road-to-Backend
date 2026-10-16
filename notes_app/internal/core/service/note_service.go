package service

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
	app_errors "github.com/Darinson16/notesgo/internal/core/errors"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type NoteService struct {
	repo output.NoteRepository
}

func NewNoteService(repo output.NoteRepository) input.NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) GetAll(ctx context.Context) (*[]schemas.GetAllNotesRow, error) {
	return s.repo.GetAll(ctx)
}

func (s *NoteService) GetAllByUserId(ctx context.Context, id int32) (*[]schemas.Note, error) {
	if id == 0 {
		return nil, app_errors.ErrIdNotProvided
	}

	return s.repo.GetAllByUserId(ctx, id)
}

func (s *NoteService) GetById(ctx context.Context, id int32) (*schemas.Note, error) {
	if id == 0 {
		return nil, app_errors.ErrIdNotProvided
	}

	return s.repo.GetById(ctx, id)
}

func (s *NoteService) Create(ctx context.Context, note *schemas.Note) (int64, error) {
	if note == nil || note.Text == "" {
		return 0, app_errors.ErrDataIsNil
	}

	return s.repo.Create(ctx, note)
}

func (s *NoteService) Update(ctx context.Context, note *schemas.Note, id int32) (*schemas.Note, error) {
	if id == 0 {
		return nil, app_errors.ErrIdNotProvided
	}

	if note == nil {
		return nil, app_errors.ErrDataIsNil
	}

	return s.repo.Update(ctx, note, id)
}

func (s *NoteService) Delete(ctx context.Context, id int32) error {
	if id == 0 {
		return app_errors.ErrIdNotProvided
	}

	return s.repo.Delete(ctx, id)
}
