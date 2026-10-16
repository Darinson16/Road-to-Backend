package input

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
)

type NoteService interface {
	GetAll(ctx context.Context) (*[]schemas.GetAllNotesRow, error)
	GetById(ctx context.Context, id int32) (*schemas.Note, error)
	Create(ctx context.Context, entity *schemas.Note) (int64, error)
	Update(ctx context.Context, entity *schemas.Note, id int32) (*schemas.Note, error)
	Delete(ctx context.Context, id int32) error
	GetAllByUserId(ctx context.Context, id int32) (*[]schemas.Note, error)
}
