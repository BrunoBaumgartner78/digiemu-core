package fs

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"digiemu-core/internal/kernel/domain"
)

type AuditByUnitReader struct {
	path string
}

func NewAuditByUnitReader(basePath string) *AuditByUnitReader {
	return &AuditByUnitReader{path: filepath.Join(basePath, "audit.ndjson")}
}

func (r *AuditByUnitReader) ListByUnitID(unitID string) ([]domain.AuditEvent, error) {
	f, err := os.Open(r.path)
	if os.IsNotExist(err) {
		return []domain.AuditEvent{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := make([]domain.AuditEvent, 0, 64)
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
		if ev.UnitID == unitID {
			out = append(out, ev)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
