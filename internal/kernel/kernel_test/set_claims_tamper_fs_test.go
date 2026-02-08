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

// Tamper detection: write claimset sidecar, modify it on disk and expect VerifyAudit to fail
func TestSetClaims_TamperDetection_FS(t *testing.T) {
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

	// set claimset
	cs := map[string]any{"schema_version": "claimset/v0", "version_id": cv.VersionID, "claims": []any{map[string]any{"id": "cl1", "text": "Original"}}}
	b, _ := json.Marshal(cs)
	setClaims := usecases.SetClaims{Repo: repo, Audit: audit, Clock: clock}
	_, err = setClaims.SetClaims(ports.SetClaimsRequest{UnitKey: "kfs", VersionID: cv.VersionID, BodyBytes: b, ActorID: "u"})
	if err != nil {
		t.Fatalf("set claims: %v", err)
	}

	// locate sidecar file and tamper it
	sidecar := filepath.Join(dir, "units", cv.UnitID+"."+cv.VersionID+".claimset.json")
	if _, err := ioutil.ReadFile(sidecar); err != nil {
		t.Fatalf("read sidecar: %v", err)
	}
	// write tampered content
	tampered := map[string]any{"schema_version": "claimset/v0", "version_id": cv.VersionID, "claims": []any{map[string]any{"id": "cl1", "text": "Tampered"}}}
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
