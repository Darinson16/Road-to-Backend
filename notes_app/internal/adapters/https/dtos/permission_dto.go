package dtos

import "github.com/Darinson16/notesgo/db/schemas"

type PermissionResponse struct {
	RoleName   string `json:"roleName"`
	ActionName string `json:"actionName"`
}

func ToPermissionResponse(p *schemas.GetAllPermissionsRow) PermissionResponse {
	return PermissionResponse{RoleName: p.Rolename, ActionName: p.Actionsname}
}
