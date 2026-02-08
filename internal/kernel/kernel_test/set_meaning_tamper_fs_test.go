package kernel_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"digiemu-core/internal/kernel/adapters/fs"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

// Tamper detection: write meaning sidecar, modify it on disk and expect VerifyAudit to fail
func TestSetMeaning_TamperDetection_FS(t *testing.T) {
	dir := t.TempDir()
	repo := fs.NewUnitRepo(dir)
	audit := fs.NewAuditLog(dir)
	reader := fs.NewAuditReader(dir)
	clock := memory.FakeClock{Now: 1700000000}

	// create unit + version
	createUnit := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	_, err := createUnit.CreateUnit(ports.CreateUnitRequest{Key: "kfs", Title: "TitleFS", ActorID: "u"})
	if err != nil {
		t.Fatalf("create unit: %v", err)
	}
	createVersion := usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}
	cv, err := createVersion.CreateVersion(ports.CreateVersionRequest{UnitKey: "kfs", Label: "v1", Content: "c", ActorID: "u"})
	if err != nil {
		t.Fatalf("create version: %v", err)
	}

	// set meaning
	m := map[string]any{"schema_version": "meaning/v1", "title": "Original"}
	b, _ := json.Marshal(m)
	setMeaning := usecases.SetMeaning{Repo: repo, Audit: audit, Clock: clock}
	_, err = setMeaning.SetMeaning(ports.SetMeaningRequest{UnitKey: "kfs", VersionID: cv.VersionID, MeaningJSON: b, ActorID: "u"})
	if err != nil {
		t.Fatalf("set meaning: %v", err)
	}

	// locate sidecar file and tamper it
	sidecar := filepath.Join(dir, "units", cv.UnitID+"."+cv.VersionID+".meaning.json")
	if _, err := ioutil.ReadFile(sidecar); err != nil {
		t.Fatalf("read sidecar: %v", err)
	}
	// write tampered content
	tampered := map[string]any{"schema_version": "meaning/v1", "title": "Tampered"}
	tb, _ := json.MarshalIndent(tampered, "", "  ")
	if err := ioutil.WriteFile(sidecar, tb, 0o644); err != nil {
		t.Fatalf("tamper write: %v", err)
	}

	// verify should detect mismatch when StrictHash=true
	verifier := usecases.VerifyAudit{Repo: repo, Audit: reader}
	out, err := verifier.VerifyAudit(ports.VerifyAuditRequest{StrictHash: true})
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if out.Ok {
		t.Fatalf("expected verify to fail due to tampering, got ok")
	}
}
