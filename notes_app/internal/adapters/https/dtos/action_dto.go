package dtos

import "github.com/Darinson16/notesgo/db/schemas"

type ActionResponse struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

func ToActionResponse(r *schemas.Action) ActionResponse {
	return ActionResponse{Id: r.ID, Name: r.Name}
}
