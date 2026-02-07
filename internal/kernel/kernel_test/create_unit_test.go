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
	uc := usecases.CreateUnit{Repo: repo}

	out, err := uc.CreateUnit(ports.CreateUnitRequest{
		Key:   "reglement-bau",
		Title: "Bau-Reglement",
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
}

func TestCreateUnit_DuplicateKey(t *testing.T) {
	repo := memory.NewUnitRepo()
	uc := usecases.CreateUnit{Repo: repo}

	_, err := uc.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title"})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}

	_, err = uc.CreateUnit(ports.CreateUnitRequest{Key: "abc", Title: "Title 2"})
	if err != domain.ErrUnitAlreadyExists {
		t.Fatalf("expected ErrUnitAlreadyExists, got %v", err)
	}
}
