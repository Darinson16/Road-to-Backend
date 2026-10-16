package repository

import (
	"context"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type AuthRepository struct {
	queries *schemas.Queries
	repo    output.UserRepository
}

func NewAuthRepository(queries *schemas.Queries, repo output.UserRepository) output.AuthRepository {
	return &AuthRepository{queries: queries, repo: repo}
}

func (r *AuthRepository) SignUp(ctx context.Context, user *schemas.User) (int64, error) {
	result, err := r.queries.AddUser(ctx, schemas.AddUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		PhoneNumber:  user.PhoneNumber,
		RoleID:       user.RoleID,
	})

	if err != nil {
		return 0, err
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return insertId, nil
}

func (r *AuthRepository) Login(ctx context.Context, email string) (*schemas.User, error) {
	return r.repo.CheckAndGetInfoFromEmail(ctx, email)
}
