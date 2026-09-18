package helper

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// AppError is the standard application error carrying an HTTP status code.
// Repository/service layers return this; the handler maps it to a JSON body.
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e == nil {
		return "unknown error"
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap supports errors.Is / errors.As.
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// NewAppError builds an AppError.
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// BadRequest creates a 400 error.
func BadRequest(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

// NotFound creates a 404 error.
func NotFound(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

// Internal creates a 500 error, wrapping the underlying cause.
func Internal(message string, err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message, Err: err}
}

// CodeOf extracts the HTTP status code from err (defaults to 500).
func CodeOf(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Code != 0 {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

// MessageOf extracts a client-safe message from err.
func MessageOf(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Message != "" {
		return appErr.Message
	}
	return http.StatusText(http.StatusInternalServerError)
}

// Kept for backward compatibility with existing callers.
type BaseErrorResponse = AppError

func FormatValidationError(err error) string {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		for _, fe := range verrs {
			switch fe.Tag() {
			case "required":
				return fe.Field() + " is required"
			case "gt":
				return fe.Field() + " must be greater than " + fe.Param()
			case "gte":
				return fe.Field() + " must be greater than or equal to " + fe.Param()
			case "max":
				return fe.Field() + " must be at most " + fe.Param() + " characters"
			case "min":
				return fe.Field() + " must be at least " + fe.Param() + " characters"
			default:
				return fe.Field() + " is invalid"
			}
		}
	}
	return "invalid request body"
}
