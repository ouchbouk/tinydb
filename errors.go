package tinydb

import "errors"

var (
	ErrKeyNotFound = errors.New("tinydb: key not found")
	ErrClosed      = errors.New("tinydb: database closed")
	ErrEmptyKey    = errors.New("tinydb: empty key")
)
