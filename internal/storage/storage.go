package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	
	ErrLinkNotFound = errors.New("link not found")
	ErrLinkExists = errors.New("link exists")
)
