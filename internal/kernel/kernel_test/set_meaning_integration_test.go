package kernel_test

import (
	"testing"

	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

// Integration test: SetMeaning -> ExportUnitSnapshot -> VerifyAudit (memory)
func TestSetMeaning_ExportAndVerify_Memory(t *testing.T) {
	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	reader := memory.NewAuditReader(audit)
	clock := memory.FakeClock{Now: 1700000000}

	// create unit + version
	createUnit := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	cu, err := createUnit.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title", ActorID: "u"})
	if err != nil {
		t.Fatalf("create unit: %v", err)
	}
	createVersion := usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}
	cv, err := createVersion.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "v1", Content: "c", ActorID: "u"})
	if err != nil {
		t.Fatalf("create version: %v", err)
	}

	// set meaning (simple payload)
	mjson := []byte(`{"schema_version":"meaning/v1","title":"M"}`)
	setMeaning := usecases.SetMeaning{Repo: repo, Audit: audit, Clock: clock}
	_, err = setMeaning.SetMeaning(ports.SetMeaningRequest{UnitKey: "abc", VersionID: cv.VersionID, MeaningJSON: mjson, ActorID: "u"})
	if err != nil {
		t.Fatalf("set meaning: %v", err)
	}

	// export snapshot and verify audit
	exporter := usecases.ExportUnitSnapshot{Repo: repo, Audit: memory.NewAuditByUnitReader(audit)}
	snap, err := exporter.ExportUnitSnapshot(ports.ExportUnitSnapshotRequest{UnitKey: "abc", IncludeAudit: true})
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if len(snap.Audit) == 0 {
		t.Fatalf("expected audit events")
	}

	verifier := usecases.VerifyAudit{Repo: repo, Audit: reader}
	out, err := verifier.VerifyAudit(ports.VerifyAuditRequest{StrictHash: true})
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !out.Ok {
		t.Fatalf("expected verify ok, got: %+v", out)
	}
	_ = cu
}
