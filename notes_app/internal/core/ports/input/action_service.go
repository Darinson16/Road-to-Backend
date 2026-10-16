package input

import (
	"github.com/Darinson16/notesgo/db/schemas"
)

type ActionService interface {
	BaseManagementService[schemas.Action]
}
