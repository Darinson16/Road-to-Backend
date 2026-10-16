package v1

import (
	"net/http"
	"strconv"

	"github.com/Darinson16/notesgo/internal/adapters/https/dtos"
	app_errors "github.com/Darinson16/notesgo/internal/core/errors"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/middlewares"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service input.RoleService
}

func NewRoleHandler(service input.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func RegisterRolesRoutes(server *gin.RouterGroup, handler *RoleHandler, userService input.UserService) {
	server.GET("/", middlewares.RBAC(userService, "Admin:Read"), handler.GetAll)
	server.GET("/:id", middlewares.RBAC(userService, "Admin:Read"), handler.GetById)
}

// GetAll godoc
// @Summary		List all roles
// @Tags			roles
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is a list of roles"
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Read"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/roles/ [get]
func (h *RoleHandler) GetAll(ctx *gin.Context) {
	roles, err := h.service.GetAll(ctx)

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the data for roles.", err.Error())
		return
	}

	rolesResponse := []dtos.RoleResponse{}
	for _, role := range *roles {
		rolesResponse = append(rolesResponse, dtos.ToRoleResponse(&role))
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Roles data fetched correctly!", rolesResponse)
}

// GetById godoc
// @Summary		Get a role by id
// @Tags			roles
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Param			id				path		int		true	"Role id"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is the role"
// @Failure		400				{object}	dtos.ResponseStruct
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Read"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/roles/{id} [get]
func (h *RoleHandler) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "There was an error parsing the provided id.", err.Error())
		return
	}

	if id == 0 {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, app_errors.ErrIdNotProvided.Error(), "There is not a valid id to search in the data")
		return
	}

	role, err := h.service.GetById(ctx, int32(id))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the role.", err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Role data fetched correctly!", dtos.ToRoleResponse(role))
}
