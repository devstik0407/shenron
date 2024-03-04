package store

import "errors"

var (
	ErrExpired  = errors.New("requested item against given key is expired")
	ErrNotFound = errors.New("requested item not found")
	ErrNoKeys   = errors.New("no keys found")
)
