package utils

import (
	"errors"
	"net/http"
)

// AppError represents a structured application error with an HTTP status code
type AppError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"error"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewError creates a custom structured AppError
func NewError(statusCode int, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// WriteError writes a structured AppError or standard error to http.ResponseWriter
func WriteError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		JSON(w, appErr.StatusCode, appErr)
		return
	}
	// Fallback for default errors
	Error(w, http.StatusInternalServerError, err.Error())
}

// Common pre-defined application errors
var (
	ErrInvalidRole        = NewError(http.StatusBadRequest, "Invalid role")
	ErrPasswordHashFailed = NewError(http.StatusInternalServerError, "Password hash failed")
	ErrUserCreationFailed = NewError(http.StatusBadRequest, "User cannot be created")
	ErrUserNotFound       = NewError(http.StatusNotFound, "User not found")
	ErrPasswordNotMatch   = NewError(http.StatusUnauthorized, "Password does not match")
	ErrTokenGenFailed     = NewError(http.StatusInternalServerError, "Token generation failed")
	ErrUnauthorized       = NewError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden          = NewError(http.StatusForbidden, "Wrong role or unauthorized")
	ErrBadRequest         = NewError(http.StatusBadRequest, "Bad request")
	ErrNotFound           = NewError(http.StatusNotFound, "Resource not found")
)
