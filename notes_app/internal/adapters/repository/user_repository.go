package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
)

type UserRepository struct {
	queries *schemas.Queries
}

func NewUserRepository(queries *schemas.Queries) output.UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) GetAll(ctx context.Context) (*[]schemas.User, error) {
	users, err := r.queries.GetAllUsers(ctx)

	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *UserRepository) GetById(ctx context.Context, id int32) (*schemas.User, error) {
	user, err := r.queries.GetUserById(ctx, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *schemas.User, id int32) (*schemas.User, error) {
	err := r.queries.UpdateUser(ctx, schemas.UpdateUserParams{
		ID:           id,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		PhoneNumber:  user.PhoneNumber,
	})

	if err != nil {
		return nil, err
	}

	updatedUser, err := r.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int32) error {
	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) ChangeRoles(ctx context.Context, userId, roleId int32) (bool, error) {
	err := r.queries.ChangeRoles(ctx, schemas.ChangeRolesParams{
		ID:     userId,
		RoleID: roleId,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) CheckAndGetInfoFromEmail(ctx context.Context, email string) (*schemas.User, error) {
	user, err := r.queries.CheckAndGetInfoFromEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("There is not rows found!. %s", err)
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) HasAccess(ctx context.Context, userId int32, permission string) (bool, error) {
	count, err := r.queries.HasAccess(ctx, schemas.HasAccessParams{
		ID:   userId,
		Name: permission,
	})

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
