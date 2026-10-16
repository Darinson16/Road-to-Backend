package middlewares

import (
	"fmt"
	"net/http"

	"github.com/Darinson16/notesgo/internal/adapters/https/dtos"
	"github.com/Darinson16/notesgo/utils"
	"github.com/gin-gonic/gin"
)

const (
	Administrator int32 = 1
	User          int32 = 2
)

func Authenticate(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		dtos.SendErrorResponse(ctx, http.StatusUnauthorized, "There is an error with the provided token", "Authorization token could not be empty")
		ctx.Abort()
		return
	}

	userEmail, userID, err := utils.VerifyToken(token)
	if err != nil {
		dtos.SendErrorResponse(ctx, http.StatusUnauthorized, "Verifying token failed!", "Provided token is not valid")
		ctx.Abort()
		return
	}

	ctx.Set("userEmail", userEmail)
	ctx.Set("userId", userID)
	ctx.Next()
}

func GetUserIdFromContext(ctx *gin.Context) (int32, error) {
	rawId, exists := ctx.Get("userId")
	if !exists {
		return 0, fmt.Errorf("There was an error getting the id.")
	}

	var ok bool
	userId, ok := rawId.(int32)
	if !ok || userId == 0 {
		return 0, fmt.Errorf("There was an error parsing the provided id.")
	}

	return userId, nil
}
