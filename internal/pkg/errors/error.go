package errors

import (
	"github.com/gin-gonic/gin"
)

// ErrorCode defines unique error codes
type ErrorCode string

const (
	ErrInvalidInput      ErrorCode = "INVALID_INPUT"
	ErrValidationFailed  ErrorCode = "VALIDATION_FAILED"
	ErrEmailNotFound     ErrorCode = "EMAIL_NOT_FOUND"
	ErrIncorrectPassword ErrorCode = "INCORRECT_PASSWORD"
	ErrTokenGeneration   ErrorCode = "TOKEN_GENERATION_FAILED"
	ErrUserNotFound      ErrorCode = "USER_NOT_FOUND"
	ErrInvalidRole       ErrorCode = "INVALID_ROLE"
	ErrCreateUserFailed  ErrorCode = "CREATE_USER_FAILED"
	ErrFetchUsersFailed  ErrorCode = "FETCH_USERS_FAILED"
	ErrUpdateUserFailed  ErrorCode = "UPDATE_USER_FAILED"
	ErrDeleteUserFailed  ErrorCode = "DELETE_USER_FAILED"
)

// APIError defines a standardized error response
type APIError struct {
	Code    int       `json:"code"`
	Error   ErrorCode `json:"error"`
	Message string    `json:"message"`
}

// HandleError sends a standardized error response
func HandleError(c *gin.Context, status int, errCode ErrorCode, message string) {
	c.JSON(status, APIError{
		Code:    status,
		Error:   errCode,
		Message: message,
	})
}
