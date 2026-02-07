package fs

import "time"

// Persistence schema for FS adapter
type VersionRecord struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type UnitRecord struct {
	ID          string          `json:"id"`
	Key         string          `json:"key"`
	Title       string          `json:"title"`
	Description string          `json:"description,omitempty"`
	CreatedAt   string          `json:"created_at"`
	Versions    []VersionRecord `json:"versions"`
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}
