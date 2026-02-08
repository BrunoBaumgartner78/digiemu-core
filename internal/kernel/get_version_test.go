package kernel_test

import (
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
	"testing"
)

func TestGetVersion_HappyPath(t *testing.T) {
	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	createUnit := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	uout, err := createUnit.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title", ActorID: "test"})
	if err != nil {
		t.Fatalf("create unit: %v", err)
	}
	_ = uout

	createVersion := usecases.CreateVersion{Repo: repo, Audit: audit, Clock: clock}
	vout, err := createVersion.CreateVersion(ports.CreateVersionRequest{UnitKey: "abc", Label: "v1", Content: "x", ActorID: "test"})
	if err != nil {
		t.Fatalf("create version: %v", err)
	}

	uc := usecases.GetVersion{Repo: repo}
	out, err := uc.GetVersion(ports.GetVersionRequest{VersionID: vout.VersionID})
	if err != nil {
		t.Fatalf("get version: %v", err)
	}
	if out.Version.ID != vout.VersionID {
		t.Fatalf("expected id %s, got %s", vout.VersionID, out.Version.ID)
	}
	if out.Version.UnitID == "" {
		t.Fatalf("expected unitID")
	}
}

func TestGetVersion_NotFound(t *testing.T) {
	repo := memory.NewUnitRepo()
	uc := usecases.GetVersion{Repo: repo}

	_, err := uc.GetVersion(ports.GetVersionRequest{VersionID: "ver_doesnotexist"})
	if err != domain.ErrVersionNotFound {
		t.Fatalf("expected ErrVersionNotFound, got %v", err)
	}
}
