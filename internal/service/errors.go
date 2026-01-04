package core

import "errors"

var (
	ErrBadRequest      = errors.New("bad request")
	ErrConflict        = errors.New("conflict")
	ErrInternal        = errors.New("internal server error")
	ErrNotFound        = errors.New("not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrUnimplemented   = errors.New("unimplemented")
	ErrAuthClientError = errors.New("auth client error")
	ErrForbidden = errors.New("forbidden")
)
