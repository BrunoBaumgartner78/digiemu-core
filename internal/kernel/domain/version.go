package domain

import "strings"

type Version struct {
	ID      string
	UnitID  string
	Label   string
	Content string

	// v0.2: lineage + auditability + integrity
	PrevVersionID string
	ContentHash   string // hex sha256
	CreatedAtUnix int64
	ActorID       string
	MeaningHash   string
	ClaimSetHash  string
}

func NewVersion(unitID, label, content string) (Version, error) {
	unitID = strings.TrimSpace(unitID)
	label = strings.TrimSpace(label)
	content = strings.TrimSpace(content)

	if unitID == "" {
		return Version{}, ErrUnitNotFound
	}
	if len(label) < 1 {
		return Version{}, ErrInvalidVersionLabel
	}
	if content == "" {
		return Version{}, ErrEmptyContent
	}

	return Version{
		ID:      NewID("ver"),
		UnitID:  unitID,
		Label:   label,
		Content: content,
	}, nil
}
