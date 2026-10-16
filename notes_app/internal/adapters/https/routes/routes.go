package routes

import (
	"net/http"

	v1 "github.com/Darinson16/notesgo/internal/adapters/https/handlers/v1"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/Darinson16/notesgo/middlewares"
	"github.com/gin-gonic/gin"
)

var containerHandlers *v1.ContainerHandler

func Register(server *gin.Engine, handlers *v1.ContainerHandler, userService input.UserService) {
	containerHandlers = handlers

	v1Group := server.Group("/api/v1")
	v1GroupAuth := server.Group("/api/v1")
	v1GroupAuth.Use(middlewares.Authenticate)
	RegisterV1Routes(v1Group, v1GroupAuth, userService)
}

func RegisterV1Routes(server, serverAuth *gin.RouterGroup, userService input.UserService) {

	server.GET("/health", HealthCheck)

	authGroup := server.Group("/auth")
	actionsGroup := serverAuth.Group("/actions")
	rolesGroup := serverAuth.Group("/roles")
	notesGroup := serverAuth.Group("/notes")
	usersGroup := serverAuth.Group("/users")
	permissionsGroup := serverAuth.Group("/permissions")

	v1.RegisterAuthRoutes(authGroup, &containerHandlers.AuthHandler)
	v1.RegisterActionsRoutes(actionsGroup, &containerHandlers.ActionHandler, userService)
	v1.RegisterUsersRoutes(usersGroup, &containerHandlers.UserHandler, userService)
	v1.RegisterNotesRoutes(notesGroup, &containerHandlers.NoteHandler, userService)
	v1.RegisterRolesRoutes(rolesGroup, &containerHandlers.RoleHandler, userService)
	v1.RegisterPermissionsRoutes(permissionsGroup, &containerHandlers.PermissionHandler, userService)
}

// HealthCheck godoc
// @Summary		Healthcheck
// @Tags			health
// @Produce		json
// @Success		200	{object}	map[string]string	"Service status"
// @Router			/health [get]
func HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Everything's good"})
}
