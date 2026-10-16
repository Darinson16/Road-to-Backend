package input

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
)

type UserService interface {
	GetAll(ctx context.Context) (*[]schemas.User, error)
	GetById(ctx context.Context, id int32) (*schemas.User, error)
	Update(ctx context.Context, entity *schemas.User, id int32) (*schemas.User, error)
	Delete(ctx context.Context, id int32) error
	HasAccess(ctx context.Context, userId int32, permission string) (bool, error)
	ChangeRoles(ctx context.Context, userId, roleId int32) (*schemas.User, error)
}
