package main

import (
	"log"

	"github.com/Darinson16/notesgo/db"
	"github.com/Darinson16/notesgo/db/schemas"
	_ "github.com/Darinson16/notesgo/docs"
	v1 "github.com/Darinson16/notesgo/internal/adapters/https/handlers/v1"
	"github.com/Darinson16/notesgo/internal/adapters/https/routes"
	"github.com/Darinson16/notesgo/internal/adapters/repository"
	"github.com/Darinson16/notesgo/internal/core/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			NotesGo API
// @version		1.0
// @description	Notes API with JWT authentication and RBAC authorization (HasAccess).
// @host			localhost:8080
// @BasePath		/api/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name							Authorization
// @description					Raw JWT from POST /auth/login (no Bearer prefix).

func main() {
	server := gin.Default()
	server.Use(gin.Recovery())

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: could not load .env file, relying on environment variables:", err)
	}

	if err := db.MigrateUp(db.MigrationDir()); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	DB := db.Connection()
	queries := schemas.New(DB)

	containerRepository := repository.NewContainerRepository(queries)
	containerService := service.NewContainerService(containerRepository)
	containerHandler := v1.NewContainerHandler(containerService)
	routes.Register(server, containerHandler, containerService.UserService)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.Run(":8080")
}
