package kernel_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

func TestUncertainty_Tamper_FS(t *testing.T) {
	dir, err := ioutil.TempDir("", "digiemu-test-uncertainty")
	if err != nil {
		t.Fatalf("tmpdir: %v", err)
	}
	defer os.RemoveAll(dir)

	repo := fsrepo.NewUnitRepo(dir)
	audit := fsrepo.NewAuditLog(dir)
	clock := memory.RealClock{}

	cu := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	_, err = cu.CreateUnit(ports.CreateUnitRequest{Key: "kfs", Title: "t", Description: "d", ActorID: "test"})
	if err != nil {
		t.Fatalf("create unit: %v", err)
	}
	cv := usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}
	outV, err := cv.CreateVersion(ports.CreateVersionRequest{UnitKey: "kfs", Label: "lbl", Content: "c", ActorID: "test"})
	if err != nil {
		t.Fatalf("create version: %v", err)
	}

	su := usecases.SetUncertainty{Repo: repo, Audit: audit, Clock: clock}
	uJSON := []byte(`{"schema_version":"uncertainty/v0","id":"u1","type":"empirical","level":"low","applies_to":{"scope":"version"}}`)
	_, err = su.SetUncertainty(ports.SetUncertaintyRequest{UnitKey: "kfs", VersionID: outV.VersionID, BodyBytes: uJSON, ActorID: "test"})
	if err != nil {
		t.Fatalf("set uncertainty: %v", err)
	}

	// tamper the sidecar
	side := filepath.Join(dir, "units", "unit_kfs."+outV.VersionID+".uncertainty.json")
	if err := ioutil.WriteFile(side, []byte("{}"), 0o644); err != nil {
		t.Fatalf("tamper write: %v", err)
	}

	ver := usecases.VerifyAudit{Repo: repo, Audit: fsrepo.NewAuditReader(dir)}
	res, err := ver.VerifyAudit(ports.VerifyAuditRequest{UnitKey: "kfs", StrictHash: true})
	if err != nil {
		t.Fatalf("verify audit: %v", err)
	}
	if res.Ok {
		t.Fatalf("expected verify audit to fail due to tamper, but ok=true")
	}
}
