package ports

import "digiemu-core/internal/kernel/domain"

type AuditLog interface {
	Append(ev domain.AuditEvent) error
}
