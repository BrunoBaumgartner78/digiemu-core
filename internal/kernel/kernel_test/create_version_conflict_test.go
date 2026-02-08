package kernel_test

import (
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
	"testing"
)

func TestCreateVersion_ConflictOnBaseVersion(t *testing.T) {
	repo := memory.NewUnitRepo()
	auditU := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	createUnit := usecases.CreateUnit{Repo: repo, Audit: auditU, Clock: clock}
	outU, err := createUnit.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title", ActorID: "test"})
	if err != nil {
		t.Fatalf("create unit: %v", err)
	}
	_ = outU

	auditV := memory.NewAuditLog()
	uc := usecases.CreateVersion{Repo: repo, Audit: auditV, Clock: clock}

	_, err = uc.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "v1", Content: "x", ActorID: "test"})
	if err != nil {
		t.Fatalf("create v1: %v", err)
	}

	_, err = uc.CreateVersion(ports.CreateVersionRequest{
		UnitKey:       "abc",
		Label:         "v2",
		Content:       "y",
		ActorID:       "test",
		BaseVersionID: "ver-wrong",
	})
	if err != domain.ErrConflict {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}
