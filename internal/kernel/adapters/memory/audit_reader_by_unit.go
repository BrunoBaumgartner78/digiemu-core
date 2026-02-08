package memory

import "digiemu-core/internal/kernel/domain"

type AuditByUnitReader struct {
	Log *AuditLog
}

func NewAuditByUnitReader(log *AuditLog) *AuditByUnitReader {
	return &AuditByUnitReader{Log: log}
}

func (r *AuditByUnitReader) ListByUnitID(unitID string) ([]domain.AuditEvent, error) {
	r.Log.mu.RLock()
	defer r.Log.mu.RUnlock()

	out := make([]domain.AuditEvent, 0, 16)
	for _, ev := range r.Log.Events {
		if ev.UnitID == unitID {
			out = append(out, ev)
		}
	}
	return out, nil
}
