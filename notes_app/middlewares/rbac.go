package middlewares

import (
	"net/http"

	"github.com/Darinson16/notesgo/internal/adapters/https/dtos"
	"github.com/Darinson16/notesgo/internal/core/ports/input"
	"github.com/gin-gonic/gin"
)

func RBAC(userService input.UserService, permissions ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawId, exists := ctx.Get("userId")
		if !exists {
			dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "User ID not found", "Authenticate must run before RBAC")
			ctx.Abort()
			return
		}

		userId, ok := rawId.(int32)
		if !ok || userId == 0 {
			dtos.SendErrorResponse(ctx, http.StatusUnauthorized, "Invalid user ID", "userId in context is not int32")
			ctx.Abort()
			return
		}

		var allowed bool
		var err error
		if len(permissions) > 1 {
			for _, permission := range permissions {
				allowed, err = userService.HasAccess(ctx.Request.Context(), userId, permission)
				if allowed {
					break
				}
			}
		} else {
			allowed, err = userService.HasAccess(ctx.Request.Context(), userId, permissions[0])
		}

		if err != nil {
			dtos.SendErrorResponse(ctx, http.StatusInternalServerError, "Could not validate permission", err.Error())
			ctx.Abort()
			return
		}

		if !allowed {
			dtos.SendErrorResponse(ctx, http.StatusForbidden, "Forbidden", "User does not have permission for this action. Communicate with your admin.")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
