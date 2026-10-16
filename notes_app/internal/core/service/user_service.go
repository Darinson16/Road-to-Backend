package service

import (
	"context"
	"fmt"

	"github.com/Darinson16/notesgo/db/schemas"
	app_errors "github.com/Darinson16/notesgo/internal/core/errors"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/internal/core/ports/output"
	"github.com/Darinson16/notesgo/utils"
)

type UserService struct {
	userRepo output.UserRepository
	roleRepo output.RoleRepository
}

func NewUserService(userRepo output.UserRepository, roleRepo output.RoleRepository) input.UserService {
	return &UserService{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *UserService) GetAll(ctx context.Context) (*[]schemas.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *UserService) GetById(ctx context.Context, id int32) (*schemas.User, error) {
	if id == 0 {
		return nil, app_errors.ErrIdNotProvided
	}

	return s.userRepo.GetById(ctx, id)
}

func (s *UserService) Update(ctx context.Context, user *schemas.User, id int32) (*schemas.User, error) {
	if id == 0 {
		return nil, app_errors.ErrIdNotProvided
	}

	if user == nil {
		return nil, app_errors.ErrDataIsNil
	}

	if _, err := s.GetById(ctx, id); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = hashedPassword

	return s.userRepo.Update(ctx, user, id)
}

func (s *UserService) Delete(ctx context.Context, id int32) error {
	if id == 0 {
		return app_errors.ErrIdNotProvided
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ChangeRoles(ctx context.Context, userId, roleId int32) (*schemas.User, error) {
	if userId == 0 || roleId == 0 {
		return nil, app_errors.ErrIdNotProvided
	}

	user, err := s.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if _, err := s.roleRepo.GetById(ctx, roleId); err != nil {
		return nil, err
	}

	_, err = s.userRepo.ChangeRoles(ctx, userId, roleId)
	if err != nil {
		return nil, err
	}

	user.RoleID = roleId
	return user, nil
}

func (s *UserService) HasAccess(ctx context.Context, userId int32, permission string) (bool, error) {
	if userId == 0 {
		return false, app_errors.ErrIdNotProvided
	}

	if permission == "" {
		return false, fmt.Errorf("There is permission to validate. must write a value.")
	}

	return s.userRepo.HasAccess(ctx, userId, permission)
}
