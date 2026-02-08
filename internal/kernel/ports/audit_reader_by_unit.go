package ports

import "digiemu-core/internal/kernel/domain"

// AuditLogByUnitReader reads audit events filtered by unitId.
// Implementations may scan full logs (FS) or filter in-memory (memory adapter).
type AuditLogByUnitReader interface {
	ListByUnitID(unitID string) ([]domain.AuditEvent, error)
}
