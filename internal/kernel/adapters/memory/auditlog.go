package memory

import (
	"digiemu-core/internal/kernel/domain"
	"sync"
)

type AuditLog struct {
	mu     sync.RWMutex
	Events []domain.AuditEvent
}

func NewAuditLog() *AuditLog {
	return &AuditLog{Events: []domain.AuditEvent{}}
}

func (l *AuditLog) Append(ev domain.AuditEvent) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Events = append(l.Events, ev)
	return nil
}

// Scan implements ports.AuditLogReader by iterating events and calling fn.
func (l *AuditLog) Scan(fn func(ev domain.AuditEvent) error) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, ev := range l.Events {
		if err := fn(ev); err != nil {
			return err
		}
	}
	return nil
}
