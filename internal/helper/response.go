package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
}

type PaginatedResponse struct {
	Response Response
	Meta     PaginatedMeta `json:"meta"`
}

type PaginatedMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

func SuccessResponse(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(ctx *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	ctx.JSON(statusCode, response)
}

func BadRequestResponse(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusBadRequest, message, err)
}

func UnauthorizedResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusUnauthorized, message, nil)
}

func ForbiddenResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusForbidden, message, nil)
}

func NotFoundResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusNotFound, message, nil)
}

func InternalServerError(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusInternalServerError, message, err)
}

func PaginatedSuccessResponse(ctx *gin.Context, message string, data any, meta PaginatedMeta) {
	ctx.JSON(http.StatusOK, PaginatedResponse{
		Response: Response{
			Success: true,
			Message: message,
			Data:    data,
		}, Meta: meta,
	})
}
