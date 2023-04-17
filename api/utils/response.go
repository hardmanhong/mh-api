package utils

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

type ApiResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func ApiSuccessResponse(data interface{}) ApiResponse {
	return ApiResponse{
		Code:    200,
		Data:    data,
		Message: "success",
	}
}
func ApiErrorResponse(code int, message string) ApiResponse {
	return ApiResponse{
		Code:    code,
		Data:    nil,
		Message: message,
	}
}
