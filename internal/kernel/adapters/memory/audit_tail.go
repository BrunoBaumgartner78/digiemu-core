package memory

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type AuditTail struct {
	Log *AuditLog
}

func NewAuditTail(log *AuditLog) *AuditTail {
	return &AuditTail{Log: log}
}

func (t *AuditTail) Tail(in ports.AuditTailRequest) ([]domain.AuditEvent, error) {
	n := in.N
	if n <= 0 {
		n = 50
	}

	t.Log.mu.RLock()
	defer t.Log.mu.RUnlock()

	out := make([]domain.AuditEvent, 0, n)
	for _, ev := range t.Log.Events {
		if in.Type != "" && ev.Type != in.Type {
			continue
		}
		if in.UnitID != "" && ev.UnitID != in.UnitID {
			continue
		}
		if in.VersionID != "" && ev.VersionID != in.VersionID {
			continue
		}
		out = append(out, ev)
		if len(out) > n {
			out = out[1:]
		}
	}
	return out, nil
}
