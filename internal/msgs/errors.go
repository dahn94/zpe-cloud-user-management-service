package msgs

import "errors"

// Define common error messages as variables.
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidRequestPayload = errors.New("invalid request payload")
	ErrForbidden             = errors.New("forbidden")
	ErrInternalServerError   = errors.New("internal server error")
	ErrUserAlreadyExists     = errors.New("user already exists")
)
