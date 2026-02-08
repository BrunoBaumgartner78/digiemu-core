package domain

// AuditEvent is append-only journal event (NDJSON friendly).
type AuditEvent struct {
	Schema  string `json:"schema"`
	ID      string `json:"id"`
	Type    string `json:"type"`
	AtUnix  int64  `json:"atUnix"`
	ActorID string `json:"actorId"`

	UnitID    string `json:"unitId,omitempty"`
	VersionID string `json:"versionId,omitempty"`

	Data any `json:"data,omitempty"`
}

type UnitCreatedData struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

type VersionCreatedData struct {
	PrevVersionID string `json:"prevVersionId,omitempty"`
	ContentHash   string `json:"contentHash"`
	Label         string `json:"label"`
}

type MeaningSetData struct {
	MeaningHash   string `json:"meaning_hash"`
	MeaningPath   string `json:"meaning_path,omitempty"`
	SchemaVersion string `json:"schema_version,omitempty"`
	InlinePreview *struct {
		Title   string `json:"title,omitempty"`
		Purpose string `json:"purpose,omitempty"`
	} `json:"inline_preview,omitempty"`
}
