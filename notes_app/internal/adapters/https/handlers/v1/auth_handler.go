package v1

import (
	"database/sql"
	"net/http"

	"github.com/Darinson16/notesgo/db/schemas"
	"github.com/Darinson16/notesgo/internal/adapters/https/dtos"
	"github.com/Darinson16/notesgo/internal/core/domain"
	app_errors "github.com/Darinson16/notesgo/internal/core/errors"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service input.AuthService
}

func NewAuthHandler(service input.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func RegisterAuthRoutes(server *gin.RouterGroup, handler *AuthHandler) {
	server.POST("/signup", handler.SignUp)
	server.POST("/login", handler.Login)
}

// SignUp godoc
// @Summary		Sign up a new user
// @Description	Creates a user with the default User role (role_id 3).
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user	body		dtos.SignUpRequest	true	"Signup payload"
// @Success		200		{object}	dtos.ResponseStruct	"Envelope; data holds the created user id"
// @Failure		400		{object}	dtos.ResponseStruct
// @Failure		500		{object}	dtos.ResponseStruct
// @Router			/auth/signup [post]
func (h *AuthHandler) SignUp(ctx *gin.Context) {
	var userData *dtos.SignUpRequest
	if err := ctx.ShouldBindJSON(&userData); err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "Could not parse body correctly", err.Error())
		return
	}

	user := &schemas.User{
		Name:         userData.Name,
		Email:        userData.Email,
		PasswordHash: userData.Password,
		PhoneNumber:  sql.NullString{String: userData.PhoneNumber, Valid: true},
		RoleID:       3,
	}

	userId, err := h.service.SignUp(ctx, user)
	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There was an error creating the user.", err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "User Created Successfully", gin.H{"id": userId})
}

// Login godoc
// @Summary		Log in and get a JWT
// @Description	Returns a raw JWT (no Bearer prefix) to send in the Authorization header.
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			credentials	body		dtos.LoginRequest	true	"Login payload"
// @Success		200			{object}	dtos.ResponseStruct	"Envelope; data holds the token"
// @Failure		400			{object}	dtos.ResponseStruct
// @Failure		401			{object}	dtos.ResponseStruct
// @Failure		500			{object}	dtos.ResponseStruct
// @Router			/auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var loginDto *dtos.LoginRequest

	if err := ctx.ShouldBindJSON(&loginDto); err != nil {
		dtos.SendErrorResponse(ctx, http.StatusBadRequest, "Could not parse body correctly. Body Malformed", err.Error())
		return
	}

	userCredentials := &domain.Login{
		Email:    loginDto.Email,
		Password: loginDto.Password,
	}

	userId, err := h.service.Login(ctx, userCredentials)
	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusUnauthorized, "Invalid user credentials! ", app_errors.ErrInvalidCredentials.Error())
		return
	}

	token, err := utils.GenerateToken(loginDto.Email, userId)
	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "There is was an error generating the token", err.Error())
		return
	}

	dtos.SendSuccessResponse(ctx, http.StatusOK, "Login successful", gin.H{"token": token})
}
