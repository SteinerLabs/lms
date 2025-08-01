package repository

import "errors"

// Common errors
var (
	// ErrNotFound is returned when a record is not found
	ErrNotFound = errors.New("record not found")

	// ErrAlreadyExists is returned when a record already exists
	ErrAlreadyExists = errors.New("record already exists")

	// ErrInvalidInput is returned when input is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrInternal is returned when an internal error occurs
	ErrInternal = errors.New("internal error")

	// ErrUnauthorized is returned when a user is not authorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when a user is forbidden from performing an action
	ErrForbidden = errors.New("forbidden")
)
