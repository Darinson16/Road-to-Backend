package output

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
)

type AuthRepository interface {
	SignUp(ctx context.Context, user *schemas.User) (int64, error)
	Login(ctx context.Context, email string) (*schemas.User, error)
}
