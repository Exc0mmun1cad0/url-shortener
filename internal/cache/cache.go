package cache

import "errors"

var (
	ErrNoPing   = errors.New("no ping")
	ErrNotFound = errors.New("key not found")
)
