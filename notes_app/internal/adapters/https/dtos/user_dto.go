package dtos

import (
	"time"

	"github.com/Darinson16/notesgo/db/schemas"
)

type UserResponse struct {
	Id          int32     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	RoleID      int32     `json:"roleId"`
}

type UpdateUserRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type ChangeRolesRequest struct {
	UserId int32 `json:"userId" binding:"required"`
	RoleId int32 `json:"roleId" binding:"required"`
}

func ToUserResponse(u *schemas.User) UserResponse {
	dto := UserResponse{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		RoleID:    u.RoleID,
	}

	if u.PhoneNumber.Valid {
		dto.PhoneNumber = u.PhoneNumber.String
	}

	return dto
}
