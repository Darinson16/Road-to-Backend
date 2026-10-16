package output

import (
	"github.com/Darinson16/notesgo/db/schemas"
)

type ActionRepository interface {
	BaseManagementRepository[schemas.Action]
}
