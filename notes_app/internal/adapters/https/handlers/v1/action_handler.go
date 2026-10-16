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

type ActionHandler struct {
	service input.ActionService
}

func NewActionHandler(service input.ActionService) *ActionHandler {
	return &ActionHandler{service: service}
}

func RegisterActionsRoutes(server *gin.RouterGroup, handler *ActionHandler, userService input.UserService) {
	server.GET("/", middlewares.RBAC(userService, "Admin:Read"), handler.GetAll)
	server.GET("/:id", middlewares.RBAC(userService, "Admin:Read"), handler.GetById)
}

// GetAll godoc
// @Summary		List all actions
// @Tags			actions
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is a list of actions"
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Read"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/actions/ [get]
func (h *ActionHandler) GetAll(ctx *gin.Context) {
	actions, err := h.service.GetAll(ctx)

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the data for actions.", err.Error())
		return
	}

	actionsResponse := []dtos.ActionResponse{}
	for _, action := range *actions {
		actionsResponse = append(actionsResponse, dtos.ToActionResponse(&action))
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "actions data fetched correctly!", actionsResponse)
}

// GetById godoc
// @Summary		Get an action by id
// @Tags			actions
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Param			id				path		int		true	"Action id"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is the action"
// @Failure		400				{object}	dtos.ResponseStruct
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Read"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/actions/{id} [get]
func (h *ActionHandler) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "There was an error parsing the provided id.", err.Error())
		return
	}

	if id == 0 {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, app_errors.ErrIdNotProvided.Error(), "There is not a valid id to search in the data")
		return
	}

	action, err := h.service.GetById(ctx, int32(id))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the action data.", err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Action data fetched correctly!", action)
}
