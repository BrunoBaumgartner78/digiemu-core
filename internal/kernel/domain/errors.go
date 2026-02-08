package domain

import "errors"

var (
	ErrInvalidUnitKey      = errors.New("invalid unit key")
	ErrInvalidUnitTitle    = errors.New("invalid unit title")
	ErrUnitAlreadyExists   = errors.New("unit already exists")
	ErrUnitNotFound        = errors.New("unit not found")
	ErrInvalidVersionLabel = errors.New("invalid version label")
	ErrEmptyContent        = errors.New("empty content")

	// uncertainty errors
	ErrInvalidSchemaVersion         = errors.New("invalid schema version")
	ErrInvalidUncertaintyID         = errors.New("invalid uncertainty id")
	ErrInvalidUncertaintyType       = errors.New("invalid uncertainty type")
	ErrInvalidUncertaintyLevel      = errors.New("invalid uncertainty level")
	ErrMissingClaimIDForUncertainty = errors.New("missing claim id for uncertainty applies_to scope=claim")
	ErrInvalidAppliesToScope        = errors.New("invalid applies_to scope")
)
