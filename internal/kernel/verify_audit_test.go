package kernel_test

import (
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
	"testing"
)

func TestVerifyAudit_Ok(t *testing.T) {
	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	createUnit := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	_, err := createUnit.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title", ActorID: "test"})
	if err != nil {
		t.Fatalf("create unit: %v", err)
	}

	createVersion := usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}
	_, err = createVersion.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "v1", Content: "x", ActorID: "test"})
	if err != nil {
		t.Fatalf("create version: %v", err)
	}

	reader := memory.NewAuditReader(audit)
	verify := usecases.VerifyAudit{Repo: repo, Audit: reader}

	out, err := verify.VerifyAudit(ports.VerifyAuditRequest{StrictHash: false})
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !out.Ok {
		t.Fatalf("expected ok, missing=%d dup=%d mismatch=%d", len(out.Missing), len(out.Duplicates), len(out.HashMismatches))
	}
}

func TestVerifyAudit_MissingDetected(t *testing.T) {
	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	// CreateUnit with audit
	createUnit := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	uout, err := createUnit.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title", ActorID: "test"})
	if err != nil {
		t.Fatalf("create unit: %v", err)
	}

	// CreateVersion WITHOUT writing to audit (simulate missing by using a separate audit log for version)
	emptyAudit := memory.NewAuditLog()
	createVersion := usecases.CreateVersion{Repo: repo, Audit: emptyAudit, Clock: clock}
	_, err = createVersion.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "v1", Content: "x", ActorID: "test"})
	if err != nil {
		t.Fatalf("create version: %v", err)
	}

	reader := memory.NewAuditReader(audit) // reader only sees unit.created, not version.created
	verify := usecases.VerifyAudit{Repo: repo, Audit: reader}

	out, err := verify.VerifyAudit(ports.VerifyAuditRequest{StrictHash: false})
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if out.Ok {
		t.Fatalf("expected not ok")
	}

	foundMissingVersion := false
	for _, m := range out.Missing {
		if m.EventType == "version.created" && m.UnitID == uout.UnitID && m.VersionID != "" {
			foundMissingVersion = true
			break
		}
	}
	if !foundMissingVersion {
		t.Fatalf("expected missing version.created")
	}
}
