package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type NoteRepository struct {
	queries *schemas.Queries
}

func NewNoteRepository(queries *schemas.Queries) output.NoteRepository {
	return &NoteRepository{queries: queries}
}

func (r *NoteRepository) GetAll(ctx context.Context) (*[]schemas.GetAllNotesRow, error) {
	notes, err := r.queries.GetAllNotes(ctx)

	if err != nil {
		return nil, err
	}

	return &notes, nil
}

func (r *NoteRepository) GetAllByUserId(ctx context.Context, id int32) (*[]schemas.Note, error) {
	notes, err := r.queries.GetAllNotesByUserId(ctx, id)

	if err != nil {
		return nil, err
	}

	return &notes, nil
}

func (r *NoteRepository) GetById(ctx context.Context, id int32) (*schemas.Note, error) {
	note, err := r.queries.GetNoteById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		return nil, err
	}

	return &note, nil
}

func (r *NoteRepository) Create(ctx context.Context, note *schemas.Note) (int64, error) {
	result, err := r.queries.AddNote(ctx, schemas.AddNoteParams{
		Text:     note.Text,
		Priority: note.Priority,
		UserID:   note.UserID,
	})

	if err != nil {
		return 0, err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return insertId, nil
}

func (r *NoteRepository) Update(ctx context.Context, note *schemas.Note, id int32) (*schemas.Note, error) {
	err := r.queries.UpdateNote(ctx, schemas.UpdateNoteParams{
		ID:       id,
		Text:     note.Text,
		Priority: note.Priority,
	})

	if err != nil {
		return nil, err
	}

	updatedNote, err := r.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return updatedNote, nil
}

func (r *NoteRepository) Delete(ctx context.Context, id int32) error {
	err := r.queries.DeleteNote(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
