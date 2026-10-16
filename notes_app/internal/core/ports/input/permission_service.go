package input

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
)

type PermissionService interface {
	GetAll(ctx context.Context) (*[]schemas.GetAllPermissionsRow, error)
}
