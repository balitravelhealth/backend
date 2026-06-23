package repository

import "errors"

var (
	ErrProviderMismatch    = errors.New("account registered with different login method")
	ErrTokenNotFound       = errors.New("token not found")
	ErrNotFound            = errors.New("record not found")
	ErrAlreadyExists       = errors.New("record already exists")
)
