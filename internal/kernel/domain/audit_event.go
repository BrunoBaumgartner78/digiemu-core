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
