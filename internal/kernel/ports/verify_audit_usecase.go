package ports

type VerifyAuditRequest struct {
	// Optional filter: verify only a single unit (by key).
	UnitKey string

	// If true, verify ContentHash matches the audit event payload for version.created.
	StrictHash bool
}

type MissingAudit struct {
	UnitID    string
	VersionID string // empty for unit.created checks
	EventType string
}

type DuplicateAudit struct {
	EventType string
	TargetID  string // UnitID or VersionID depending on EventType
}

type HashMismatch struct {
	UnitID       string
	VersionID    string
	ExpectedHash string
	EventHash    string
}

type VerifyAuditResponse struct {
	TotalUnits    int
	TotalVersions int

	Missing        []MissingAudit
	Duplicates     []DuplicateAudit
	HashMismatches []HashMismatch

	Ok bool
}

type VerifyAuditUsecase interface {
	VerifyAudit(in VerifyAuditRequest) (VerifyAuditResponse, error)
}
