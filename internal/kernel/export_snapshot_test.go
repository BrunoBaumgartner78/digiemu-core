package kernel_test

import (
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
	"testing"
)

func TestExportUnitSnapshot_OrderAndAudit(t *testing.T) {
	repo := memory.NewUnitRepo()
	auditLog := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	createUnit := usecases.CreateUnit{Repo: repo, Audit: auditLog, Clock: clock}
	_, err := createUnit.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title", ActorID: "test"})
	if err != nil {
		t.Fatalf("create unit: %v", err)
	}

	createVersion := usecases.CreateVersion{Repo: repo, Audit: auditLog, Clock: clock}
	_, err = createVersion.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "v1", Content: "one", ActorID: "test"})
	if err != nil {
		t.Fatalf("create v1: %v", err)
	}
	_, err = createVersion.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "v2", Content: "two", ActorID: "test"})
	if err != nil {
		t.Fatalf("create v2: %v", err)
	}

	ucNoAudit := usecases.ExportUnitSnapshot{Repo: repo, Audit: nil}
	snap, err := ucNoAudit.ExportUnitSnapshot(ports.ExportUnitSnapshotRequest{UnitKey: "abc", IncludeAudit: false})
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if len(snap.Versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(snap.Versions))
	}
	if snap.Versions[0].Label != "v1" || snap.Versions[1].Label != "v2" {
		t.Fatalf("expected order v1,v2 got %s,%s", snap.Versions[0].Label, snap.Versions[1].Label)
	}
	if len(snap.Audit) != 0 {
		t.Fatalf("expected no audit in snapshot")
	}

	reader := memory.NewAuditByUnitReader(auditLog)
	ucAudit := usecases.ExportUnitSnapshot{Repo: repo, Audit: reader}
	snap2, err := ucAudit.ExportUnitSnapshot(ports.ExportUnitSnapshotRequest{UnitKey: "abc", IncludeAudit: true})
	if err != nil {
		t.Fatalf("export with audit: %v", err)
	}
	// unit.created + two version.created = 3
	if len(snap2.Audit) != 3 {
		t.Fatalf("expected 3 audit events, got %d", len(snap2.Audit))
	}
}
