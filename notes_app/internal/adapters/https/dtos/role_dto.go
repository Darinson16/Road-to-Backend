package dtos

import "github.com/Darinson16/notesgo/db/schemas"

type RoleResponse struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func ToRoleResponse(r *schemas.Role) RoleResponse {
	return RoleResponse{Id: r.ID, Name: r.Name}
}
