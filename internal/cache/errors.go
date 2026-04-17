package cache

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrInternal             = errors.New("internal error")
	ErrInvalidCacheInstance = errors.New("invalid cache instance")
)
