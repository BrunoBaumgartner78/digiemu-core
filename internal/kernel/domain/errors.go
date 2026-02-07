package domain

import "errors"

var (
	ErrInvalidUnitKey      = errors.New("invalid unit key")
	ErrInvalidUnitTitle    = errors.New("invalid unit title")
	ErrUnitAlreadyExists   = errors.New("unit already exists")
	ErrUnitNotFound        = errors.New("unit not found")
	ErrInvalidVersionLabel = errors.New("invalid version label")
	ErrEmptyContent        = errors.New("empty content")
)
