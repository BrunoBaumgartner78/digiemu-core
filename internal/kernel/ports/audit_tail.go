package ports

import "digiemu-core/internal/kernel/domain"

type AuditTailRequest struct {
	// Last N events (default 50 if 0)
	N int

	// Optional filters
	Type      string
	UnitID    string
	VersionID string
}

// AuditTailReader can return the last N events (best-effort) without scanning full log in CLI.
type AuditTailReader interface {
	Tail(in AuditTailRequest) ([]domain.AuditEvent, error)
}
