package registration

import (
	"errors"
	"fmt"
)

// Sentinel errors (predefined errors to compare against)
var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
	ErrInvalidEmail = errors.New("invalid email format")
	ErrInvalidUser  = errors.New("invalid user")
	ErrInvalidAge   = errors.New("invalid age")
)

// Custom error type with more context
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s': %s", e.Field, e.Message)
}
