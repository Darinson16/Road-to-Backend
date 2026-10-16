package output

import (
	"github.com/Darinson16/notesgo/db/schemas"
)

type RoleRepository interface {
	BaseManagementRepository[schemas.Role]
}
