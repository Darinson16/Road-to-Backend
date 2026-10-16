package service

import (
	"github.com/Darinson16/notesgo/internal/adapters/repository"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
)

type ContainerService struct {
	AuthService       input.AuthService
	ActionService     input.ActionService
	PermissionService input.PermissionService
	RoleService       input.RoleService
	UserService       input.UserService
	NoteService       input.NoteService
}

func NewContainerService(repos *repository.ContainerRepository) *ContainerService {
	return &ContainerService{
		ActionService:     NewActionService(repos.ActionRepository),
		PermissionService: NewPermissionService(repos.PermissionRepository),
		RoleService:       NewRoleService(repos.RoleRepository),
		UserService:       NewUserService(repos.UserRepository, repos.RoleRepository),
		NoteService:       NewNoteService(repos.NoteRepository),
		AuthService:       NewAuthService(repos.AuthRepository),
	}
}
