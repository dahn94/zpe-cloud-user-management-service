package msgs

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrInvalidRequestPayload   = errors.New("invalid request payload")
	ErrForbidden               = errors.New("forbidden")
	ErrInternalServerError     = errors.New("internal server error")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrInvalidRole             = errors.New("invalid role")
	ErrInsufficientPermissions = errors.New("insufficient permissions to assign role")
	ErrMethodNotAllowed        = errors.New("method not allowed")
)
