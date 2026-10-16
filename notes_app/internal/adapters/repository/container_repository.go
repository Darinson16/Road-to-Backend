package repository

import (
	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type ContainerRepository struct {
	AuthRepository       output.AuthRepository
	ActionRepository     output.ActionRepository
	PermissionRepository output.PermissionRepository
	RoleRepository       output.RoleRepository
	UserRepository       output.UserRepository
	NoteRepository       output.NoteRepository
}

func NewContainerRepository(queries *schemas.Queries) *ContainerRepository {
	return &ContainerRepository{
		ActionRepository:     NewActionRepository(queries),
		PermissionRepository: NewPermissionRepository(queries),
		RoleRepository:       NewRoleRepository(queries),
		UserRepository:       NewUserRepository(queries),
		NoteRepository:       NewNoteRepository(queries),
		AuthRepository:       NewAuthRepository(queries, NewUserRepository(queries)),
	}
}
