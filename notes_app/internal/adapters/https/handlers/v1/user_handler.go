package v1

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/adapters/https/dtos"
	app_errors "github.com/Darinson16/notesgo/internal/core/errors"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/middlewares"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service input.UserService
}

func NewUserHandler(service input.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func RegisterUsersRoutes(server *gin.RouterGroup, handler *UserHandler, userService input.UserService) {
	server.GET("/", middlewares.RBAC(userService, "Admin:Read", "Manager:Read"), handler.GetAll)
	server.GET("/:id", middlewares.RBAC(userService, "Admin:Read", "Manager:Read"), handler.GetById)
	server.POST("/change-role", middlewares.RBAC(userService, "Admin:Update"), handler.ChangeRoles)
	server.PUT("/:id", middlewares.RBAC(userService, "Admin:Update", "Manager:Update"), handler.Update)
	server.DELETE("/:id", middlewares.RBAC(userService, "Admin:Delete"), handler.Delete)
}

// GetAll godoc
// @Summary		List all users
// @Tags			users
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is a list of users"
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Read or Manager:Read"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/users/ [get]
func (h *UserHandler) GetAll(ctx *gin.Context) {
	users, err := h.service.GetAll(ctx)

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the data for users.", err.Error())
		return
	}

	userDto := []dtos.UserResponse{}
	for _, user := range *users {
		userDto = append(userDto, dtos.ToUserResponse(&user))
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Users fetched correctly!", userDto)
}

// GetById godoc
// @Summary		Get a user by id
// @Tags			users
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Param			id				path		int		true	"User id"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is the user"
// @Failure		400				{object}	dtos.ResponseStruct
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Read or Manager:Read"
// @Failure		404				{object}	dtos.ResponseStruct
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/users/{id} [get]
func (h *UserHandler) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "There was an error parsing the provided id.", err.Error())
		return
	}

	if id == 0 {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, app_errors.ErrIdNotProvided.Error(), "There is not a valid id to search data")
		return
	}

	user, err := h.service.GetById(ctx, int32(id))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the user.", err.Error())
		return
	}

	if user == nil {
		dtos.SendErrorResponse(ctx, http.StatusNotFound, "There is not register user with this id.", "User not found")
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "User data fetched correctly!", dtos.ToUserResponse(user))
}

// Update godoc
// @Summary		Update a user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string					true	"Raw JWT from /auth/login"
// @Param			id				path		int						true	"User id"
// @Param			user			body		dtos.UpdateUserRequest	true	"User payload"
// @Success		200				{object}	dtos.ResponseStruct		"Envelope; data is the updated user"
// @Failure		400				{object}	dtos.ResponseStruct
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Update or Manager:Update"
// @Failure		404				{object}	dtos.ResponseStruct
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/users/{id} [put]
func (h *UserHandler) Update(ctx *gin.Context) {
	var userDto *dtos.UpdateUserRequest
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "There was an error parsing the provided id.", err.Error())
		return
	}

	if id == 0 {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, app_errors.ErrIdNotProvided.Error(), "There is not a valid id to search in the data")
		return
	}

	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "Could not parse body correctly", err.Error())
		return
	}

	user := &schemas.User{
		Name:         userDto.Name,
		Email:        userDto.Email,
		PasswordHash: userDto.Password,
		PhoneNumber:  sql.NullString{String: userDto.PhoneNumber},
	}

	updatedUser, err := h.service.Update(ctx, user, int32(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			dtos.SendErrorResponse(ctx, http.StatusNotFound, "There is not register user with this id.", err.Error())
			return
		}
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error updating the user.", err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "User Updated Successfully", dtos.ToUserResponse(updatedUser))
}

// Delete godoc
// @Summary		Delete a user
// @Tags			users
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Param			id				path		int		true	"User id"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is the deleted id"
// @Failure		400				{object}	dtos.ResponseStruct
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Delete"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/users/{id} [delete]
func (h *UserHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "There was an error parsing the provided id.", err.Error())
		return
	}

	if id == 0 {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, app_errors.ErrIdNotProvided.Error(), "There is not a valid id to search in the data")
		return
	}

	if err = h.service.Delete(ctx, int32(id)); err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the user.", err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "User deleted!", id)
}

// ChangeRoles godoc
// @Summary		Change a user's role
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string						true	"Raw JWT from /auth/login"
// @Param			roles			body		dtos.ChangeRolesRequest	true	"User and role ids"
// @Success		200				{object}	dtos.ResponseStruct			"Envelope; data is the updated user"
// @Failure		400				{object}	dtos.ResponseStruct
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Update"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/users/change-role [post]
func (h *UserHandler) ChangeRoles(ctx *gin.Context) {
	var changesRolesDto *dtos.ChangeRolesRequest

	if err := ctx.ShouldBindJSON(&changesRolesDto); err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "There was an error parsing the body.", err.Error())
		return
	}

	user, err := h.service.ChangeRoles(ctx, changesRolesDto.UserId, changesRolesDto.RoleId)
	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error changing the role. Try later!", err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Role changed successfully", dtos.ToUserResponse(user))
}
