package kernel_test

import (
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
	"testing"
)

func TestExportUnitSnapshot_HashDeterministic(t *testing.T) {
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

	uc := usecases.ExportUnitSnapshot{Repo: repo, Audit: memory.NewAuditByUnitReader(auditLog)}

	s1, err := uc.ExportUnitSnapshot(ports.ExportUnitSnapshotRequest{UnitKey: "abc", IncludeAudit: false})
	if err != nil {
		t.Fatalf("export1: %v", err)
	}
	s2, err := uc.ExportUnitSnapshot(ports.ExportUnitSnapshotRequest{UnitKey: "abc", IncludeAudit: false})
	if err != nil {
		t.Fatalf("export2: %v", err)
	}
	if s1.SnapshotHash == "" {
		t.Fatalf("expected snapshotHash")
	}
	if s1.SnapshotHash != s2.SnapshotHash {
		t.Fatalf("expected deterministic snapshotHash, got %s != %s", s1.SnapshotHash, s2.SnapshotHash)
	}
	if s1.AuditHash != "" {
		t.Fatalf("expected empty auditHash when IncludeAudit=false")
	}

	s3, err := uc.ExportUnitSnapshot(ports.ExportUnitSnapshotRequest{UnitKey: "abc", IncludeAudit: true})
	if err != nil {
		t.Fatalf("export3: %v", err)
	}
	if s3.AuditHash == "" {
		t.Fatalf("expected auditHash when IncludeAudit=true")
	}
}

func TestExportUnitSnapshot_HashChangesOnContent(t *testing.T) {
	repo := memory.NewUnitRepo()
	auditLog := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	createUnit := usecases.CreateUnit{Repo: repo, Audit: auditLog, Clock: clock}
	_, _ = createUnit.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title", ActorID: "test"})

	createVersion := usecases.CreateVersion{Repo: repo, Audit: auditLog, Clock: clock}
	_, _ = createVersion.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "v1", Content: "one", ActorID: "test"})

	uc := usecases.ExportUnitSnapshot{Repo: repo, Audit: memory.NewAuditByUnitReader(auditLog)}
	before, err := uc.ExportUnitSnapshot(ports.ExportUnitSnapshotRequest{UnitKey: "abc", IncludeAudit: false})
	if err != nil {
		t.Fatalf("export before: %v", err)
	}

	// add second version with different content
	_, _ = createVersion.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "v2", Content: "DIFFERENT", ActorID: "test"})

	after, err := uc.ExportUnitSnapshot(ports.ExportUnitSnapshotRequest{UnitKey: "abc", IncludeAudit: false})
	if err != nil {
		t.Fatalf("export after: %v", err)
	}

	if before.SnapshotHash == after.SnapshotHash {
		t.Fatalf("expected snapshotHash to change when versions change")
	}
}
