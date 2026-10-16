package v1

import (
	"github.com/Darinson16/notesgo/internal/core/service"
)

type ContainerHandler struct {
	AuthHandler       AuthHandler
	ActionHandler     ActionHandler
	PermissionHandler PermissionHandler
	RoleHandler       RoleHandler
	UserHandler       UserHandler
	NoteHandler       NoteHandler
}

func NewContainerHandler(services *service.ContainerService) *ContainerHandler {
	return &ContainerHandler{
		ActionHandler:     *NewActionHandler(services.ActionService),
		PermissionHandler: *NewPermissionHandler(services.PermissionService),
		RoleHandler:       *NewRoleHandler(services.RoleService),
		UserHandler:       *NewUserHandler(services.UserService),
		NoteHandler:       *NewNoteHandler(services.NoteService),
		AuthHandler:       *NewAuthHandler(services.AuthService),
	}
}
