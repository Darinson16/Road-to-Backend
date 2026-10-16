package output

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
)

type PermissionRepository interface {
	GetAll(ctx context.Context) (*[]schemas.GetAllPermissionsRow, error)
}
