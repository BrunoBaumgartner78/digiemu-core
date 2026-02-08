package kernel_test

import (
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
	"testing"
)

func TestCreateUnit_HappyPath(t *testing.T) {
	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	uc := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}

	out, err := uc.CreateUnit(ports.CreateUnitRequest{
		Key:     "reglement-bau",
		Title:   "Bau-Reglement",
		ActorID: "test",
	})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if out.UnitID == "" {
		t.Fatalf("expected id")
	}
	if out.Key != "reglement-bau" {
		t.Fatalf("unexpected key: %s", out.Key)
	}
	if len(audit.Events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(audit.Events))
	}
	if audit.Events[0].Type != "unit.created" {
		t.Fatalf("unexpected audit type: %s", audit.Events[0].Type)
	}
}

func TestCreateUnit_DuplicateKey(t *testing.T) {
	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	uc := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}

	_, err := uc.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title", ActorID: "test"})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}

	_, err = uc.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title 2", ActorID: "test"})
	if err != domain.ErrUnitAlreadyExists {
		t.Fatalf("expected ErrUnitAlreadyExists, got %v", err)
	}
}
