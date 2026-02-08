package memory

import "digiemu-core/internal/kernel/domain"

type AuditReader struct {
	Log *AuditLog
}

func NewAuditReader(log *AuditLog) *AuditReader {
	return &AuditReader{Log: log}
}

func (r *AuditReader) Scan(fn func(ev domain.AuditEvent) error) error {
	r.Log.mu.RLock()
	defer r.Log.mu.RUnlock()
	for _, ev := range r.Log.Events {
		if err := fn(ev); err != nil {
			return err
		}
	}
	return nil
}
