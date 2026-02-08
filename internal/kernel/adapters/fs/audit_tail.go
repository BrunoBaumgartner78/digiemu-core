package fs

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type AuditTail struct {
	path string
}

func NewAuditTail(basePath string) *AuditTail {
	return &AuditTail{path: filepath.Join(basePath, "audit.ndjson")}
}

func (t *AuditTail) Tail(in ports.AuditTailRequest) ([]domain.AuditEvent, error) {
	n := in.N
	if n <= 0 {
		n = 50
	}

	f, err := os.Open(t.path)
	if os.IsNotExist(err) {
		return []domain.AuditEvent{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Simple ring buffer via slice
	buf := make([]domain.AuditEvent, 0, n)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev domain.AuditEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			return nil, err
		}
		if !matchAuditTailFilters(ev, in) {
			continue
		}
		if len(buf) < n {
			buf = append(buf, ev)
		} else {
			// shift left (O(n) per overflow), acceptable for small n
			copy(buf, buf[1:])
			buf[len(buf)-1] = ev
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return buf, nil
}

func matchAuditTailFilters(ev domain.AuditEvent, in ports.AuditTailRequest) bool {
	if in.Type != "" && ev.Type != in.Type {
		return false
	}
	if in.UnitID != "" && ev.UnitID != in.UnitID {
		return false
	}
	if in.VersionID != "" && ev.VersionID != in.VersionID {
		return false
	}
	return true
}
