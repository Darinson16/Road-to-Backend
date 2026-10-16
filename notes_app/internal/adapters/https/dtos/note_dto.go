package dtos

import (
	"time"

	"github.com/Darinson16/notesgo/db/schemas"
)

type GetAllNoteResponse struct {
	Id        int32     `json:"id"`
	Text      string    `json:"text"`
	Priority  int16     `json:"priority"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserInfo  struct {
		Id    int32  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"userInfo"`
}

type NoteResponse struct {
	Id        int32     `json:"id"`
	Text      string    `json:"text"`
	Priority  int16     `json:"priority"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserId    int32     `json:"userId"`
}

type NoteRequest struct {
	Text     string `json:"text" binding:"required"`
	Priority int16  `json:"priority" binding:"required"`
}

func ToGetAllNoteResponse(n *schemas.GetAllNotesRow) GetAllNoteResponse {
	return GetAllNoteResponse{
		Id:        n.ID,
		Text:      n.Text,
		Priority:  n.Priority.Int16,
		UpdatedAt: n.UpdatedAt,
		UserInfo: struct {
			Id    int32  "json:\"id\""
			Name  string "json:\"name\""
			Email string "json:\"email\""
		}{
			Id:    n.ID_2,
			Name:  n.Name,
			Email: n.Email,
		},
	}
}
func ToNoteResponse(n *schemas.Note) NoteResponse {
	return NoteResponse{
		Id:        n.ID,
		Text:      n.Text,
		Priority:  n.Priority.Int16,
		UpdatedAt: n.UpdatedAt,
		UserId:    n.UserID,
	}
}
