package fs

import (
	"encoding/json"
	"os"
	"path/filepath"

	"digiemu-core/internal/kernel/domain"
)

// AuditLog is an append-only NDJSON log.
type AuditLog struct {
	path string
}

func NewAuditLog(basePath string) *AuditLog {
	return &AuditLog{path: filepath.Join(basePath, "audit.ndjson")}
}

func (l *AuditLog) Append(ev domain.AuditEvent) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(append(b, '\n'))
	return err
}
