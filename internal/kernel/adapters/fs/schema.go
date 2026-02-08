package fs

import "time"

// Persistence schema for FS adapter
type VersionRecord struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`

	// v0.2
	PrevVersionID string `json:"prev_version_id,omitempty"`
	ContentHash   string `json:"content_hash,omitempty"`
	ActorID       string `json:"actor_id,omitempty"`
}

type UnitRecord struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`

	// v0.2
	HeadVersionID string          `json:"head_version_id,omitempty"`
	Versions      []VersionRecord `json:"versions"`
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}
