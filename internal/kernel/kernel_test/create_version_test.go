package kernel_test

import (
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
	"testing"
)

func TestCreateVersion_HappyPath(t *testing.T) {
	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	createUnit := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	unitOut, err := createUnit.CreateUnit(ports.CreateUnitRequest{
		Key:     "merkblatt-steuern",
		Title:   "Merkblatt Steuern",
		ActorID: "test",
	})
	if err != nil {
		t.Fatalf("create unit err: %v", err)
	}

	vcAudit := memory.NewAuditLog()
	uc := usecases.CreateVersion{Repo: repo, Audit: vcAudit, Clock: clock}

	out, err := uc.CreateVersion(ports.CreateVersionRequest{
		UnitKey: "merkblatt-steuern",
		Label:   "v1",
		Content: "Inhalt 1",
		ActorID: "test",
	})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if out.UnitID != unitOut.UnitID {
		t.Fatalf("expected unitID match")
	}

	u, ok, _ := repo.FindUnitByID(unitOut.UnitID)
	if !ok {
		t.Fatalf("expected unit exists")
	}
	if u.HeadVersionID == "" {
		t.Fatalf("expected headVersionID set")
	}

	vs, _ := repo.ListVersionsByUnitID(unitOut.UnitID)
	if len(vs) != 1 {
		t.Fatalf("expected 1 version, got %d", len(vs))
	}
	if len(vcAudit.Events) != 1 {
		t.Fatalf("expected 1 version audit event, got %d", len(vcAudit.Events))
	}
	if vcAudit.Events[0].Type != "version.created" {
		t.Fatalf("unexpected audit type: %s", vcAudit.Events[0].Type)
	}
}

func TestCreateVersion_UnitNotFound(t *testing.T) {
	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	uc := usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}

	_, err := uc.CreateVersion(ports.CreateVersionRequest{
		UnitKey: "does-not-exist",
		Label:   "v1",
		Content: "x",
		ActorID: "test",
	})
	if err != domain.ErrUnitNotFound {
		t.Fatalf("expected ErrUnitNotFound, got %v", err)
	}
}

func TestCreateVersion_Validation(t *testing.T) {
	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	createUnit := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	_, _ = createUnit.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title", ActorID: "test"})

	uc := usecases.CreateVersion{Repo: repo, Audit: memory.NewAuditLog(), Clock: clock}

	_, err := uc.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "", Content: "x", ActorID: "test"})
	if err != domain.ErrInvalidVersionLabel {
		t.Fatalf("expected ErrInvalidVersionLabel, got %v", err)
	}

	_, err = uc.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "v1", Content: "", ActorID: "test"})
	if err != domain.ErrEmptyContent {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
}
