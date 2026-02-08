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
