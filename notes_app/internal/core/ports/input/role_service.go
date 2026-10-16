package input

import "github.com/Darinson16/notesgo/db/schemas"

type RoleService interface {
	BaseManagementService[schemas.Role]
}
