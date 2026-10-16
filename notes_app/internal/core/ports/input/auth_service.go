package input

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/core/domain"
)

type AuthService interface {
	SignUp(ctx context.Context, user *schemas.User) (int64, error)
	Login(ctx context.Context, userCredentials *domain.Login) (int32, error)
}
