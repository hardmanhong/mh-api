package utils

import "github.com/gin-gonic/gin"

func FormatResult(code int, data interface{}, message string) gin.H {
	return gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	}
}
