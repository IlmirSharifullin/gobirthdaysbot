package common

import (
	"errors"
)

var (
	ErrNoUser            = errors.New("update's user is nil")
	ErrNotDate           = errors.New("given string is not a date")
	ErrBirthdayTypeError = errors.New("birthday type is not *storage.Birthday")
)
