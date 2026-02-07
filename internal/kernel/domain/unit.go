package domain

import "strings"

type Unit struct {
	ID    string
	Key   string
	Title string
}

func NewUnit(key, title string) (Unit, error) {
	key = strings.TrimSpace(key)
	title = strings.TrimSpace(title)

	if len(key) < 3 {
		return Unit{}, ErrInvalidUnitKey
	}
	if len(title) < 3 {
		return Unit{}, ErrInvalidUnitTitle
	}

	return Unit{
		ID:    NewID("unit"),
		Key:   key,
		Title: title,
	}, nil
}
