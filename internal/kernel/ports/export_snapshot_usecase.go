package ports

import "digiemu-core/internal/kernel/domain"

type ExportUnitSnapshotRequest struct {
	UnitKey      string
	IncludeAudit bool
}

// ExportUnitSnapshotResponse is a stable snapshot of a unit.
// Versions MUST be ordered oldest -> newest.
// Audit events (if included) are returned in file order (append order).
type ExportUnitSnapshotResponse struct {
	Unit     UnitDTO             `json:"unit"`
	Versions []VersionDTO        `json:"versions"`
	Audit    []domain.AuditEvent `json:"audit,omitempty"`

	// v0.2.5: deterministic hashes for signing/archiving
	SnapshotHash string `json:"snapshotHash"`
	AuditHash    string `json:"auditHash,omitempty"`
}

type ExportUnitSnapshotUsecase interface {
	ExportUnitSnapshot(in ExportUnitSnapshotRequest) (ExportUnitSnapshotResponse, error)
}
