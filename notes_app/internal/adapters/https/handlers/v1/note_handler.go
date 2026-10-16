package v1

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/adapters/https/dtos"
	app_errors "github.com/Darinson16/notesgo/internal/core/errors"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/middlewares"
	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	service input.NoteService
}

func NewNoteHandler(service input.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

func RegisterNotesRoutes(server *gin.RouterGroup, handler *NoteHandler, userService input.UserService) {
	server.GET("/", middlewares.RBAC(userService, "Admin:Read", "Manager:Read"), handler.GetAll)
	server.GET("/:id", middlewares.RBAC(userService, "Admin:Read", "Manager:Read"), handler.GetById)
	server.GET("/user", middlewares.RBAC(userService, "Admin:Read", "Manager:Read", "User:Read"), handler.GetAllByUserId)
	server.POST("/", middlewares.RBAC(userService, "Admin:Add", "Manager:Add", "User:Add"), handler.Create)
	server.PUT("/:id", middlewares.RBAC(userService, "Admin:Update", "Manager:Update"), handler.Update)
	server.DELETE("/:id", middlewares.RBAC(userService, "Admin:Delete"), handler.Delete)
}

// GetAll godoc
// @Summary		List all notes
// @Tags			notes
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is a list of notes with user info"
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Read or Manager:Read"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/notes/ [get]
func (h *NoteHandler) GetAll(ctx *gin.Context) {
	notes, err := h.service.GetAll(ctx)

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the data for notes.", err.Error())
		return
	}

	notesResponse := []dtos.GetAllNoteResponse{}
	for _, note := range *notes {
		notesResponse = append(notesResponse, dtos.ToGetAllNoteResponse(&note))
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Notes fetched correctly!", notesResponse)
}

// GetAllByUserId godoc
// @Summary		List notes of the authenticated user
// @Description	User id is taken from the JWT context, not from the URL.
// @Tags			notes
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is a list of notes"
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Read, Manager:Read or User:Read"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/notes/user [get]
func (h *NoteHandler) GetAllByUserId(ctx *gin.Context) {
	userId, err := middlewares.GetUserIdFromContext(ctx)
	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "userId not found!", err.Error())
		return
	}

	notes, err := h.service.GetAllByUserId(ctx, int32(userId))
	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the note.", err.Error())
		return
	}

	notesResponse := []dtos.NoteResponse{}
	for _, note := range *notes {
		notesResponse = append(notesResponse, dtos.ToNoteResponse(&note))
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Note data fetched correctly!", notesResponse)
}

// GetById godoc
// @Summary		Get a note by id
// @Tags			notes
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Param			id				path		int		true	"Note id"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is the note"
// @Failure		400				{object}	dtos.ResponseStruct
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Read or Manager:Read"
// @Failure		404				{object}	dtos.ResponseStruct
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/notes/{id} [get]
func (h *NoteHandler) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "There was an error parsing the provided id.", err.Error())
		return
	}

	if id == 0 {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, app_errors.ErrIdNotProvided.Error(), "There is not a valid id to search in the data")
	}

	note, err := h.service.GetById(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			dtos.SendErrorResponse(ctx, http.StatusNotFound, "There was not note with this id.", err.Error())
			return
		}

		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the note", err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Note data fetched correctly!", dtos.ToNoteResponse(note))
}

// Create godoc
// @Summary		Create a note for the authenticated user
// @Description	Owner is taken from the JWT context.
// @Tags			notes
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string				true	"Raw JWT from /auth/login"
// @Param			note			body		dtos.NoteRequest	true	"Note payload"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data holds the created note id"
// @Failure		400				{object}	dtos.ResponseStruct
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Add, Manager:Add or User:Add"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/notes/ [post]
func (h *NoteHandler) Create(ctx *gin.Context) {
	var noteRequest *dtos.NoteRequest
	if err := ctx.ShouldBindJSON(&noteRequest); err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "Could not parse body correctly", err.Error())
		return
	}

	userId, err := middlewares.GetUserIdFromContext(ctx)
	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "userId not found!", err.Error())
		return
	}

	note := &schemas.Note{
		Text:     noteRequest.Text,
		Priority: sql.NullInt16{Int16: noteRequest.Priority, Valid: true},
		UserID:   userId,
	}

	noteId, err := h.service.Create(ctx, note)

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("There was an error creating the note. userId: %v", userId), err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Note Created Successfully", gin.H{"id": noteId})
}

// Update godoc
// @Summary		Update a note
// @Tags			notes
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string				true	"Raw JWT from /auth/login"
// @Param			id				path		int					true	"Note id"
// @Param			note			body		dtos.NoteRequest	true	"Note payload"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is the updated note"
// @Failure		400				{object}	dtos.ResponseStruct
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Update or Manager:Update"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/notes/{id} [put]
func (h *NoteHandler) Update(ctx *gin.Context) {
	var noteUpdate *dtos.NoteRequest
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "There was an error parsing the provided id.", err.Error())
		return
	}

	if id == 0 {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, app_errors.ErrIdNotProvided.Error(), "There is not a valid id to search in the data")
	}

	if err := ctx.ShouldBindJSON(&noteUpdate); err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "Could not parse body correctly", err.Error())
		return
	}

	note := &schemas.Note{
		Text:     noteUpdate.Text,
		Priority: sql.NullInt16{Int16: noteUpdate.Priority, Valid: true},
	}

	updatedNote, err := h.service.Update(ctx, note, int32(id))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error updating the note.", err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Note Updated Successfully", dtos.ToNoteResponse(updatedNote))
}

// Delete godoc
// @Summary		Delete a note
// @Tags			notes
// @Produce		json
// @Param			Authorization	header		string	true	"Raw JWT from /auth/login"
// @Param			id				path		int		true	"Note id"
// @Success		200				{object}	dtos.ResponseStruct	"Envelope; data is the deleted id"
// @Failure		400				{object}	dtos.ResponseStruct
// @Failure		401				{object}	dtos.ResponseStruct
// @Failure		403				{object}	dtos.ResponseStruct	"Missing Admin:Delete"
// @Failure		500				{object}	dtos.ResponseStruct
// @Security		ApiKeyAuth
// @Router			/notes/{id} [delete]
func (h *NoteHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "There was an error parsing the provided id.", err.Error())
		return
	}

	if id == 0 {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, app_errors.ErrIdNotProvided.Error(), errors.New("There is not a valid id to search in the data").Error())
	}

	if err = h.service.Delete(ctx, int32(id)); err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error getting the note.", err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Note deleted!", id)
}
