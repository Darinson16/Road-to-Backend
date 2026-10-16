package dtos

import (
	"github.com/gin-gonic/gin"
)

type ResponseStruct struct {
	Message string         `json:"message"`
	Success bool           `json:"success"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Data    any            `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Details string `json:"details"`
}

func SendSuccessResponse(ctx *gin.Context, status int, message string, data any) {

	response := ResponseStruct{
		Success: true,
		Message: message,
		Data:    data,
	}

	ctx.JSON(status, response)
}

func SendErrorResponse(ctx *gin.Context, status int, message, details string) {
	response := ResponseStruct{
		Success: false,
		Message: message,
		Error: &ErrorResponse{
			Code:    status,
			Details: details,
		},
	}

	ctx.JSON(status, response)
}
