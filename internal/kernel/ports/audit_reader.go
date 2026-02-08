package ports

import "digiemu-core/internal/kernel/domain"

// AuditLogReader provides streaming access to the append-only audit log.
// Scan must stop early if fn returns an error.
type AuditLogReader interface {
	Scan(fn func(ev domain.AuditEvent) error) error
}
