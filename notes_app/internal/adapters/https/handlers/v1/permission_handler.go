package v1

import (
	"net/http"

	"github.com/Darinson16/notesgo/internal/adapters/https/dtos"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/middlewares"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	service input.PermissionService
}

func NewPermissionHandler(service input.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

func RegisterPermissionsRoutes(server *gin.RouterGroup, handler *PermissionHandler, userService input.UserService) {
	server.GET("/", middlewares.RBAC(userService, "Admin:Read"), handler.GetAll)
}

// GetAll godoc
// @Summary		List role-action permissions
// @Tags			permissions
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is a list of role/action pairs"
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Read"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/permissions/ [get]
func (h *PermissionHandler) GetAll(ctx *gin.Context) {
	permissions, err := h.service.GetAll(ctx.Request.Context())

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an internal error bringing the data", err.Error())
		return
	}

	permissionsResponse := []dtos.PermissionResponse{}
	for _, permission := range *permissions {
		permissionsResponse = append(permissionsResponse, dtos.ToPermissionResponse(&permission))
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Permissions fetch correctly!", permissionsResponse)
}
