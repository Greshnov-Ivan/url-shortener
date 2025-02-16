package serverrors

import "errors"

var (
	ErrURLNotFound = errors.New("URL not found")
	ErrURLExpired  = errors.New("URL expired")
	ErrExpired     = errors.New("date cannot be expired")
)
