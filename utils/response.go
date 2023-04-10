package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrNotFound = errors.New("not found")

func ApiSuccessResponse(data interface{}) gin.H {
	return gin.H{
		"code":    200,
		"data":    data,
		"message": "success",
	}
}
func ApiErrorResponse(code int, message string) gin.H {
	return gin.H{
		"code":    code,
		"data":    nil,
		"message": message,
	}
}
