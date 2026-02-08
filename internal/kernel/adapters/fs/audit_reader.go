package fs

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"digiemu-core/internal/kernel/domain"
)

// AuditReader scans the append-only NDJSON audit log.
type AuditReader struct {
	path string
}

func NewAuditReader(basePath string) *AuditReader {
	return &AuditReader{path: filepath.Join(basePath, "audit.ndjson")}
}

func (r *AuditReader) Scan(fn func(ev domain.AuditEvent) error) error {
	f, err := os.Open(r.path)
	if os.IsNotExist(err) {
		// no audit log yet: treat as empty
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev domain.AuditEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			return err
		}
		if err := fn(ev); err != nil {
			return err
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}
